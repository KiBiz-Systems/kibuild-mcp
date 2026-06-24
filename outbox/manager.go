package outbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/priyabratasahoo21/kibuild-mcp/validator"
)

// ResolveOutboxPath resolves the project outbox path using project.json configuration or default folder structure.
// If project.json is not found at projectPath, it checks the parent directory so that paths like
// "project/files" correctly resolve to "project/Outbox" rather than "project/files/Outbox".
func ResolveOutboxPath(projectPath string) string {
	outboxDir := "Outbox"
	// Try projectPath, then its parent — handles cases where projectPath is a files/ subfolder.
	candidates := []string{projectPath}
	if parent := filepath.Dir(projectPath); parent != projectPath && parent != "" {
		candidates = append(candidates, parent)
	}
	for _, candidate := range candidates {
		manifestPath := filepath.Join(candidate, "project.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			continue
		}
		var manifest struct {
			Structure map[string]string `json:"structure"`
		}
		if err := json.Unmarshal(data, &manifest); err == nil && manifest.Structure != nil {
			if val, exists := manifest.Structure["outbox"]; exists && val != "" {
				outboxDir = val
			}
		}
		return filepath.Join(candidate, outboxDir)
	}
	return filepath.Join(projectPath, outboxDir)
}

// LoadManifest reads and parses the manifest.json file from the outbox path.
func LoadManifest(projectPath string) (*Manifest, error) {
	outboxPath := ResolveOutboxPath(projectPath)
	manifestPath := filepath.Join(outboxPath, "manifest.json")

	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		return &Manifest{
			Artifacts: make(map[string]*Artifact),
		}, nil
	}

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read outbox manifest: %w", err)
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse outbox manifest: %w", err)
	}
	if m.Artifacts == nil {
		m.Artifacts = make(map[string]*Artifact)
	}
	return &m, nil
}

// SaveManifest serializes and writes the manifest to the project's outbox path.
func SaveManifest(projectPath string, m *Manifest) error {
	outboxPath := ResolveOutboxPath(projectPath)
	if err := os.MkdirAll(outboxPath, 0755); err != nil {
		return fmt.Errorf("failed to create outbox directory: %w", err)
	}

	manifestPath := filepath.Join(outboxPath, "manifest.json")
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal outbox manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write outbox manifest: %w", err)
	}
	return nil
}

// CreateArtifact registers a new generated artifact in the outbox manifest.
func CreateArtifact(projectPath, id, artType, name, database string) (*Artifact, error) {
	m, err := LoadManifest(projectPath)
	if err != nil {
		return nil, err
	}

	// Clean/slugify the artifact ID if needed
	id = strings.TrimSpace(strings.ToLower(id))
	id = strings.ReplaceAll(id, " ", "_")

	if art, exists := m.Artifacts[id]; exists {
		return art, nil
	}

	art := &Artifact{
		ID:        id,
		Type:      artType,
		Name:      name,
		Database:  database,
		Status:    StatusDraft,
		Versions:  []Version{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m.Artifacts[id] = art
	if err := SaveManifest(projectPath, m); err != nil {
		return nil, err
	}
	return art, nil
}

var stepNumRegex = regexp.MustCompile(`(?i)^(\s*)step\s*\d+\s*[:.)\-–—]\s*`)

func stripStepNumbers(content string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = stepNumRegex.ReplaceAllString(line, "${1}")
	}
	return strings.Join(lines, "\n")
}

// WriteVersion creates a new version directory and saves the corresponding generated files.
func WriteVersion(projectPath, id string, files map[string]string) (*Version, error) {
	m, err := LoadManifest(projectPath)
	if err != nil {
		return nil, err
	}

	art, exists := m.Artifacts[id]
	if !exists {
		return nil, fmt.Errorf("artifact %q not found in outbox manifest", id)
	}

	versionID := GenerateVersionID(len(art.Versions))
	outboxPath := ResolveOutboxPath(projectPath)
	outboxDir, errRel := filepath.Rel(projectPath, outboxPath)
	if errRel != nil || outboxDir == "." {
		outboxDir = "Outbox"
	}
	
	// Create paths: e.g. Outbox/scripts/create_contact/Version 1/
	typeFolder := art.Type + "s"
	versionRelDir := filepath.Join(outboxDir, typeFolder, art.ID, versionID)
	versionFullDir := filepath.Join(projectPath, versionRelDir)

	if err := os.MkdirAll(versionFullDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create version directory: %w", err)
	}

	// Normalize keys for script/layout artifacts to ensure suffixes match
	normalizedFiles := make(map[string]string)
	for filename, content := range files {
		lowKey := strings.ToLower(filename)
		if art.Type == "script" || art.Type == "layout" {
			if lowKey == "xml_content" || lowKey == "xml" {
				normalizedFiles[art.ID+".xml"] = content
				continue
			}
			if lowKey == "txt_content" || lowKey == "txt" {
				normalizedFiles[art.ID+".txt"] = content
				continue
			}
		}
		normalizedFiles[filename] = content
	}

	var writtenFiles []string
	for filename, content := range normalizedFiles {
		// Filter files: for script and layout artifacts, only write .xml and .txt files
		if (art.Type == "script" || art.Type == "layout") &&
			!strings.HasSuffix(strings.ToLower(filename), ".xml") &&
			!strings.HasSuffix(strings.ToLower(filename), ".txt") {
			continue
		}

		// Also sanitize step numbers for script.txt and layout.txt
		if filename == "script.txt" || filename == "layout.txt" || strings.HasSuffix(strings.ToLower(filename), ".txt") {
			content = stripStepNumbers(content)
		}

		// Name files: {Script Name}_Version N.ext  — readable in sidebar
		ext := filepath.Ext(filename)
		targetFilename := filename
		if art.Type == "script" || art.Type == "layout" {
			targetFilename = fmt.Sprintf("%s_%s%s", art.Name, versionID, ext)
		}

		filePath := filepath.Join(versionFullDir, targetFilename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return nil, fmt.Errorf("failed to write version file %q: %w", targetFilename, err)
		}
		// Record relative path in files array
		writtenFiles = append(writtenFiles, filepath.Join(versionRelDir, targetFilename))

		// Maintain "latest" version directly under artifact slug directory: Outbox/scripts/create_contact/create_contact_latest.xml
		latestName := fmt.Sprintf("%s_latest%s", art.ID, ext)
		latestPath := filepath.Join(outboxPath, typeFolder, art.ID, latestName)
		_ = os.WriteFile(latestPath, []byte(content), 0644)
	}

	// Validate XML files in the version and write validation.json
	var combinedResult = validator.ValidationResult{Passed: true}
	for filename, content := range files {
		if strings.HasSuffix(strings.ToLower(filename), ".xml") {
			vres, errVal := validator.ValidateXML(content, projectPath, art.Database)
			if errVal == nil {
				if !vres.Passed {
					combinedResult.Passed = false
				}
				// prefix message with filename for clarity
				for _, err := range vres.Errors {
					err.Message = fmt.Sprintf("[%s] %s", filename, err.Message)
					combinedResult.Errors = append(combinedResult.Errors, err)
				}
				for _, warn := range vres.Warnings {
					warn.Message = fmt.Sprintf("[%s] %s", filename, warn.Message)
					combinedResult.Warnings = append(combinedResult.Warnings, warn)
				}
			}
		}
	}
	valPath := filepath.Join(versionFullDir, "validation.json")
	valData, errValMarshal := json.MarshalIndent(combinedResult, "", "  ")
	if errValMarshal == nil {
		_ = os.WriteFile(valPath, valData, 0644)
		writtenFiles = append(writtenFiles, filepath.Join(versionRelDir, "validation.json"))
	}

	v := Version{
		VersionID: versionID,
		Timestamp: time.Now(),
		Files:     writtenFiles,
		Status:    StatusDraft,
	}

	art.Versions = append(art.Versions, v)
	art.LatestVer = versionID
	art.Status = StatusDraft
	art.UpdatedAt = time.Now()

	if err := SaveManifest(projectPath, m); err != nil {
		return nil, err
	}
	return &v, nil
}

// UpdateStatus sets the review status of an outbox artifact version.
func UpdateStatus(projectPath, id, versionID, status string, override bool) error {
	m, err := LoadManifest(projectPath)
	if err != nil {
		return err
	}

	art, exists := m.Artifacts[id]
	if !exists {
		return fmt.Errorf("artifact %q not found in outbox manifest", id)
	}

	newStatus := ArtifactStatus(status)
	if newStatus != StatusDraft && newStatus != StatusAccepted && newStatus != StatusRejected && newStatus != StatusApplied {
		return fmt.Errorf("invalid status value: %q", status)
	}

	// Validate before moving status to accepted or applied
	if newStatus == StatusAccepted || newStatus == StatusApplied {
		var versionRelDir string
		for _, v := range art.Versions {
			if v.VersionID == versionID {
				for _, f := range v.Files {
					if filepath.Base(f) == "validation.json" {
						versionRelDir = filepath.Dir(f)
						break
					}
				}
				break
			}
		}
		if versionRelDir != "" {
			valJsonPath := filepath.Join(projectPath, versionRelDir, "validation.json")
			if valData, errRead := os.ReadFile(valJsonPath); errRead == nil {
				var valRes struct {
					Passed bool `json:"passed"`
				}
				if json.Unmarshal(valData, &valRes) == nil && !valRes.Passed && !override {
					return fmt.Errorf("validation failed for this version and override is false")
				}
			}
		}
	}

	// Update specific version status
	versionFound := false
	for i, v := range art.Versions {
		if v.VersionID == versionID {
			art.Versions[i].Status = newStatus
			versionFound = true
			break
		}
	}
	if !versionFound {
		return fmt.Errorf("version %q not found for artifact %q", versionID, id)
	}

	// If updating the latest version, propagate status to the artifact itself
	if art.LatestVer == versionID {
		art.Status = newStatus
	}
	art.UpdatedAt = time.Now()

	// If status is accepted/applied, copy files to latest version folder
	if newStatus == StatusAccepted || newStatus == StatusApplied {
		outboxPath := ResolveOutboxPath(projectPath)
		typeFolder := art.Type + "s"

		for _, v := range art.Versions {
			if v.VersionID == versionID {
				for _, relFile := range v.Files {
					filename := filepath.Base(relFile)
					srcPath := filepath.Join(projectPath, relFile)
					ext := filepath.Ext(filename)
					latestName := fmt.Sprintf("%s_latest%s", art.ID, ext)
					destPath := filepath.Join(outboxPath, typeFolder, art.ID, latestName)
					
					if data, errRead := os.ReadFile(srcPath); errRead == nil {
						_ = os.WriteFile(destPath, data, 0644)
					}
				}
			}
		}
	}

	return SaveManifest(projectPath, m)
}

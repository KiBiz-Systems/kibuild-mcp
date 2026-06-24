package skills

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed default_skills/*.md
var defaultSkillsFS embed.FS

// InitDefaultSkills extracts the embedded default skills to target directory if not present.
func InitDefaultSkills(skillsDir string) error {
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return err
	}

	// Sync from workspace source if found
	if wsDir := findWorkspaceSkills("default_skills"); wsDir != "" {
		if entries, err := os.ReadDir(wsDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
					continue
				}
				srcPath := filepath.Join(wsDir, entry.Name())
				destPath := filepath.Join(skillsDir, entry.Name())
				if data, err := os.ReadFile(srcPath); err == nil {
					_ = os.WriteFile(destPath, data, 0644)
				}
			}
		}
	}

	// Load from embedded FS
	entries, err := fs.ReadDir(defaultSkillsFS, "default_skills")
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		destPath := filepath.Join(skillsDir, entry.Name())
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			data, err := defaultSkillsFS.ReadFile("default_skills/" + entry.Name())
			if err != nil {
				return err
			}
			if err := os.WriteFile(destPath, data, 0644); err != nil {
				return err
			}
		}
	}
	return nil
}

func findWorkspaceSkills(subDirName string) string {
	if execPath, err := os.Executable(); err == nil {
		dir := filepath.Dir(execPath)
		for i := 0; i < 5; i++ {
			checkPath := filepath.Join(dir, "sidecar", "skills", subDirName)
			if info, err := os.Stat(checkPath); err == nil && info.IsDir() {
				return checkPath
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	if home, err := os.UserHomeDir(); err == nil {
		checkPath := filepath.Join(home, "Documents", "KiBuild Plugin", "sidecar", "skills", subDirName)
		if info, err := os.Stat(checkPath); err == nil && info.IsDir() {
			return checkPath
		}
	}
	return ""
}

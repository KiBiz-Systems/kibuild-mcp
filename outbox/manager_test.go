package outbox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestOutboxWorkflow(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kibuild-outbox-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 1. Resolve Outbox Path
	outboxPath := ResolveOutboxPath(tempDir)
	expectedPath := filepath.Join(tempDir, "Outbox")
	if outboxPath != expectedPath {
		t.Errorf("ResolveOutboxPath = %q; want %q", outboxPath, expectedPath)
	}

	// 2. Load Manifest (should be empty first)
	manifest, err := LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	if len(manifest.Artifacts) != 0 {
		t.Errorf("expected 0 artifacts, got %d", len(manifest.Artifacts))
	}

	// 3. Create Artifact
	art, err := CreateArtifact(tempDir, "Test Script", "script", "Test Script Name", "TestDB")
	if err != nil {
		t.Fatalf("CreateArtifact failed: %v", err)
	}
	if art.ID != "test_script" {
		t.Errorf("expected artifact ID 'test_script', got %q", art.ID)
	}
	if art.Status != StatusDraft {
		t.Errorf("expected artifact status 'draft', got %q", art.Status)
	}

	// Verify it saved to manifest.json
	manifest, err = LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("LoadManifest reload failed: %v", err)
	}
	if _, exists := manifest.Artifacts["test_script"]; !exists {
		t.Fatal("expected 'test_script' to exist in loaded manifest")
	}

	// 4. Write Version
	files := map[string]string{
		"script.xml": "<fmxmlsnippet>Test Script snippet</fmxmlsnippet>",
		"readme.txt": "Instructions for running test script.",
	}
	version, err := WriteVersion(tempDir, "test_script", files)
	if err != nil {
		t.Fatalf("WriteVersion failed: %v", err)
	}
	if version.VersionID != "Version 1" {
		t.Errorf("expected version ID 'Version 1', got %q", version.VersionID)
	}
	if len(version.Files) != 3 {
		t.Errorf("expected 3 version files, got %d", len(version.Files))
	}

	// Verify version folder contains the files named: {art.Name}_Version 1.ext
	verFullDir := filepath.Join(outboxPath, "scripts", "test_script", version.VersionID)
	for name := range files {
		if !strings.HasSuffix(name, ".xml") && !strings.HasSuffix(name, ".txt") {
			continue
		}
		ext := filepath.Ext(name)
		targetName := fmt.Sprintf("Test Script Name_%s%s", version.VersionID, ext)
		fullPath := filepath.Join(verFullDir, targetName)
		if _, err := os.ReadFile(fullPath); err != nil {
			t.Errorf("failed to read version file %q: %v", targetName, err)
		}
	}

	// Verify "latest" files are generated
	latestXmlPath := filepath.Join(outboxPath, "scripts", "test_script", "test_script_latest.xml")
	data, err := os.ReadFile(latestXmlPath)
	if err != nil {
		t.Errorf("failed to read latestXmlPath: %v", err)
	}
	if string(data) != files["script.xml"] {
		t.Errorf("content of latest script.xml = %q; want %q", string(data), files["script.xml"])
	}

	// 5. Update Status
	err = UpdateStatus(tempDir, "test_script", version.VersionID, "accepted", false)
	if err != nil {
		t.Fatalf("UpdateStatus failed: %v", err)
	}

	// Verify updated status in reload
	manifest, err = LoadManifest(tempDir)
	if err != nil {
		t.Fatalf("LoadManifest reload failed: %v", err)
	}
	updatedArt := manifest.Artifacts["test_script"]
	if updatedArt.Status != StatusAccepted {
		t.Errorf("expected artifact status 'accepted', got %q", updatedArt.Status)
	}
	if updatedArt.Versions[0].Status != StatusAccepted {
		t.Errorf("expected version status 'accepted', got %q", updatedArt.Versions[0].Status)
	}
}

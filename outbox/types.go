package outbox

import "time"

type ArtifactStatus string

const (
	StatusDraft    ArtifactStatus = "draft"
	StatusAccepted ArtifactStatus = "accepted"
	StatusRejected ArtifactStatus = "rejected"
	StatusApplied  ArtifactStatus = "applied"
)

type Version struct {
	VersionID string         `json:"version_id"` // e.g. "v001_2026-05-29_1701"
	Timestamp time.Time      `json:"timestamp"`
	Files     []string       `json:"files"`      // Relative paths inside the version folder
	Status    ArtifactStatus `json:"status"`
}

type Artifact struct {
	ID        string         `json:"id"`   // slug, e.g. "create_contact"
	Type      string         `json:"type"` // "script", "layout", "schema", "doc"
	Name      string         `json:"name"`
	Database  string         `json:"database"`
	LatestVer string         `json:"latest_version"`
	Status    ArtifactStatus `json:"status"`
	Versions  []Version      `json:"versions"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Manifest struct {
	Artifacts map[string]*Artifact `json:"artifacts"`
}

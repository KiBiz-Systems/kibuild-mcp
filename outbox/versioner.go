package outbox

import (
	"fmt"
)

// GenerateVersionID returns a human-readable version label (e.g. "Version 1", "Version 2").
func GenerateVersionID(existingCount int) string {
	nextNum := existingCount + 1
	return fmt.Sprintf("Version %d", nextNum)
}

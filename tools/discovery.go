package tools

import (
	"fmt"
	"strings"
	"sync"

	"kibuild/mcp/skills"
)

var (
	globalSkillsDir string
	skillsDirMu     sync.RWMutex
)

// SetSkillsDir registers the skills directory so load_skill can use it.
func SetSkillsDir(dir string) {
	skillsDirMu.Lock()
	defer skillsDirMu.Unlock()
	globalSkillsDir = dir
}

func getSkillsDir() string {
	skillsDirMu.RLock()
	defer skillsDirMu.RUnlock()
	return globalSkillsDir
}

// LoadSkillTool returns the full markdown content of a skill by ID.
func LoadSkillTool(skillID string) (string, error) {
	if skillID == "" {
		return "", fmt.Errorf("skill_id is required")
	}

	dir := getSkillsDir()
	if dir == "" {
		return "", fmt.Errorf("skills directory not configured")
	}

	allSkills, err := skills.ListSkills(dir)
	if err != nil {
		return "", fmt.Errorf("failed to list skills: %w", err)
	}

	for _, s := range allSkills {
		if s.ID == skillID {
			if s.Content == "" {
				return "", fmt.Errorf("skill %q exists but has no content", skillID)
			}
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("# Skill: %s\n\n", s.Name))
			sb.WriteString(s.Content)
			return sb.String(), nil
		}
	}

	var available []string
	for _, s := range allSkills {
		available = append(available, s.ID)
	}
	return "", fmt.Errorf("skill %q not found. Available skills: %s", skillID, strings.Join(available, ", "))
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/priyabratasahoo21/kibuild-mcp/guides"
)

// aiGuide describes one guide file to install: the source name inside
// guides.FS, the destination function that resolves the target path, and a
// human-readable label for the setup output.
type aiGuide struct {
	src   string
	dst   func() (string, error)
	label string
}

// allGuides is the authoritative list of every guide file --setup installs.
var allGuides = []aiGuide{
	{
		src:   "claude_init.md",
		dst:   claudeCommandPath("init-kibuild-project.md"),
		label: "~/.claude/commands/init-kibuild-project.md  (Claude Code — /init-kibuild-project)",
	},
	{
		src:   "claude_setup.md",
		dst:   claudeCommandPath("setup-kibuild.md"),
		label: "~/.claude/commands/setup-kibuild.md         (Claude Code — /setup-kibuild)",
	},
	{
		src:   "codex_global.md",
		dst:   codexInstructionsPath,
		label: "~/.codex/instructions.md                    (OpenAI Codex CLI)",
	},
	{
		src:   "gemini_global.md",
		dst:   geminiGuidePath,
		label: "~/.gemini/GEMINI.md                         (Google Gemini CLI)",
	},
	{
		src:   "cursor_global.mdc",
		dst:   cursorRulesPath,
		label: "~/.cursor/rules/kibuild.mdc                 (Cursor)",
	},
}

// installAIGuides writes all guide files to their respective global locations.
// Called as the final step of runSetup().
func installAIGuides() {
	fmt.Println("[5/5] Installing AI tool guide files...")

	ok := 0
	for _, g := range allGuides {
		data, err := guides.FS.ReadFile(g.src)
		if err != nil {
			fmt.Printf("      ✗ %s\n        (could not read embedded file: %v)\n", g.label, err)
			continue
		}
		dest, err := g.dst()
		if err != nil {
			fmt.Printf("      ✗ %s\n        (could not resolve path: %v)\n", g.label, err)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			fmt.Printf("      ✗ %s\n        (could not create directory: %v)\n", g.label, err)
			continue
		}
		if err := os.WriteFile(dest, data, 0644); err != nil {
			fmt.Printf("      ✗ %s\n        (write failed: %v)\n", g.label, err)
			continue
		}
		fmt.Printf("      ✓ %s\n", g.label)
		ok++
	}

	fmt.Println()
	if ok == len(allGuides) {
		fmt.Println("      All AI guide files installed.")
	} else {
		fmt.Printf("      %d/%d files installed. See errors above.\n", ok, len(allGuides))
	}
	fmt.Println()
	fmt.Println("      Next: open any FileMaker project folder and run")
	fmt.Println("        /init-kibuild-project  (Claude Code)")
	fmt.Println("      to generate project-level guide files for your team.")
}

// ── Path resolvers ───────────────────────────────────────────────────────────

// claudeCommandPath returns a resolver for ~/.claude/commands/<name>.
func claudeCommandPath(name string) func() (string, error) {
	return func() (string, error) {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".claude", "commands", name), nil
	}
}

// codexInstructionsPath resolves ~/.codex/instructions.md (OpenAI Codex CLI).
func codexInstructionsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".codex", "instructions.md"), nil
}

// geminiGuidePath resolves ~/.gemini/GEMINI.md (Google Gemini CLI).
func geminiGuidePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gemini", "GEMINI.md"), nil
}

// cursorRulesPath resolves the Cursor global rules directory.
// Cursor reads user-level rules from ~/.cursor/rules/ on all platforms.
func cursorRulesPath() (string, error) {
	if runtime.GOOS == "windows" {
		appdata := os.Getenv("APPDATA")
		if appdata == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appdata = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appdata, "Cursor", "User", "rules", "kibuild.mdc"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cursor", "rules", "kibuild.mdc"), nil
}

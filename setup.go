package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/priyabratasahoo21/kibuild-mcp/config"
	"github.com/priyabratasahoo21/kibuild-mcp/tools"
)

const (
	githubRepo   = "priyabratasahoo21/kibuild-mcp"
	latestAPIURL = "https://api.github.com/repos/" + githubRepo + "/releases/latest"
	releaseBase  = "https://github.com/" + githubRepo + "/releases/latest/download/"
	docsURL      = "https://github.com/" + githubRepo
)

// runSetup is the interactive `kibuild-mcp --setup` flow. It self-checks the
// version against GitHub, optionally self-updates the binary, writes the Claude
// Code MCP config, and verifies that every tool is accessible — all natively,
// so it behaves identically on macOS, Linux, and Windows.
func runSetup() int {
	in := bufio.NewReader(os.Stdin)

	fmt.Println("┌──────────────────────────────────────────────┐")
	fmt.Println("│  KiBuild MCP — Setup                          │")
	fmt.Println("└──────────────────────────────────────────────┘")
	fmt.Printf("Installed version: %s\n", Version)
	fmt.Printf("Platform:          %s/%s\n\n", runtime.GOOS, runtime.GOARCH)

	// ── Step 1: Version check + optional self-update ────────────────────────
	fmt.Println("[1/5] Checking for the latest version...")
	latest, err := fetchLatestVersion()
	if err != nil {
		fmt.Printf("      Could not reach GitHub (%v).\n", err)
		fmt.Println("      Skipping update check — continuing with the installed binary.")
	} else if isNewer(latest, Version) {
		fmt.Printf("      A newer version is available: %s (you have %s)\n", latest, Version)
		fmt.Print("      Download and replace this binary now? [Y/n]: ")
		ans := strings.ToLower(strings.TrimSpace(readLine(in)))
		if ans == "" || ans == "y" || ans == "yes" {
			if err := selfUpdate(); err != nil {
				fmt.Printf("      Self-update failed: %v\n", err)
				fmt.Println("      Re-run the install script to update manually:")
				fmt.Println("        curl -fsSL https://raw.githubusercontent.com/" + githubRepo + "/main/install.sh | sh")
			} else {
				fmt.Printf("      ✓ Updated to %s. Re-run `kibuild-mcp --setup` to finish.\n", latest)
				return 0
			}
		} else {
			fmt.Println("      Keeping the installed binary.")
		}
	} else {
		fmt.Printf("      ✓ You are on the latest version (%s).\n", Version)
	}
	fmt.Println()

	// ── Step 2: Project path ────────────────────────────────────────────────
	fmt.Println("[2/5] FileMaker project folder")
	fmt.Println("      The folder that contains (or will contain) files/Schema/.")
	def := suggestProjectPath()
	if def != "" {
		fmt.Printf("      Project path [%s]: ", def)
	} else {
		fmt.Print("      Project path: ")
	}
	projectPath := strings.TrimSpace(readLine(in))
	if projectPath == "" {
		projectPath = def
	}
	if projectPath == "" {
		fmt.Println("      No path given — you can re-run --setup later. Continuing without it.")
	} else {
		if abs, err := filepath.Abs(projectPath); err == nil {
			projectPath = abs
		}
		if fi, err := os.Stat(projectPath); err != nil || !fi.IsDir() {
			fmt.Printf("      Note: %s does not exist yet — it will be used as configured.\n", projectPath)
		}
	}
	fmt.Println()

	// ── Step 3: Write MCP config ────────────────────────────────────────────
	fmt.Println("[3/5] Writing Claude Code MCP config...")
	exePath, err := os.Executable()
	if err != nil || exePath == "" {
		exePath = defaultBinaryPath()
	} else if resolved, rerr := filepath.EvalSymlinks(exePath); rerr == nil {
		exePath = resolved
	}
	cfgPath, werr := writeClaudeConfig(exePath, projectPath)
	if werr != nil {
		fmt.Printf("      ✗ Failed to write config: %v\n", werr)
	} else {
		fmt.Printf("      ✓ %s\n", cfgPath)
		fmt.Printf("        command = %s\n", exePath)
		if projectPath != "" {
			fmt.Printf("        KIBUILD_ACTIVE_PROJECT = %s\n", projectPath)
		}
	}
	fmt.Println()

	// ── Step 4: Verify tools ────────────────────────────────────────────────
	fmt.Println("[4/5] Verifying tools are accessible...")
	verifyTools()
	fmt.Println()

	// ── Step 5: Install AI guide files ──────────────────────────────────────
	installAIGuides()

	fmt.Println("────────────────────────────────────────────────")
	fmt.Println("  Restart Claude Code (close and reopen), then")
	fmt.Println("  run /mcp to confirm 'kibuild' is connected.")
	fmt.Println()
	fmt.Println("  Docs: " + docsURL)
	fmt.Println("────────────────────────────────────────────────")
	return 0
}

func readLine(r *bufio.Reader) string {
	line, _ := r.ReadString('\n')
	return strings.TrimRight(line, "\r\n")
}

// suggestProjectPath proposes a sensible default for the project folder.
func suggestProjectPath() string {
	if p := strings.TrimSpace(os.Getenv("KIBUILD_ACTIVE_PROJECT")); p != "" {
		return p
	}
	if home, err := os.UserHomeDir(); err == nil {
		if data, err := os.ReadFile(filepath.Join(home, ".fm_ai_bridge", "active_project.txt")); err == nil {
			if p := strings.TrimSpace(string(data)); p != "" {
				return p
			}
		}
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return ""
}

// defaultBinaryPath returns the conventional install location per OS.
func defaultBinaryPath() string {
	if runtime.GOOS == "windows" {
		if la := os.Getenv("LOCALAPPDATA"); la != "" {
			return filepath.Join(la, "Programs", "kibuild-mcp", "kibuild-mcp.exe")
		}
		return "kibuild-mcp.exe"
	}
	return "/usr/local/bin/kibuild-mcp"
}

// claudeConfigPath returns the path to ~/.claude.json.
func claudeConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".claude.json"), nil
}

// writeClaudeConfig safely merges the kibuild MCP server entry into the existing
// Claude Code config, preserving every other field and server.
func writeClaudeConfig(exePath, projectPath string) (string, error) {
	path, err := claudeConfigPath()
	if err != nil {
		return "", err
	}

	cfg := map[string]interface{}{}
	if data, err := os.ReadFile(path); err == nil && len(data) > 0 {
		if uerr := json.Unmarshal(data, &cfg); uerr != nil {
			// Don't clobber an unparseable config; back it up first.
			backup := path + ".bak"
			_ = os.WriteFile(backup, data, 0644)
			fmt.Printf("      Existing config was not valid JSON — backed up to %s\n", backup)
			cfg = map[string]interface{}{}
		}
	}

	servers, ok := cfg["mcpServers"].(map[string]interface{})
	if !ok || servers == nil {
		servers = map[string]interface{}{}
		cfg["mcpServers"] = servers
	}

	env := map[string]interface{}{}
	if projectPath != "" {
		env["KIBUILD_ACTIVE_PROJECT"] = projectPath
	}
	servers["kibuild"] = map[string]interface{}{
		"command": exePath,
		"env":     env,
	}

	out, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return "", err
	}
	out = append(out, '\n')
	if err := os.WriteFile(path, out, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// verifyTools prints the exact set of tools the MCP server will expose, using
// the real registry + safety filter so the report is authoritative.
func verifyTools() {
	// Mirror main()'s disabled-tool check.
	if m, err := config.NewManager(""); err == nil {
		cfgManager = m
	}

	all := tools.GetToolsSchema()
	var visible []string
	for _, t := range all {
		if tools.IsSafeMCPTool(t.Name) && !isToolDisabled(t.Name) {
			visible = append(visible, t.Name)
		}
	}

	fmt.Printf("      ✓ %d tools will be exposed to MCP clients.\n", len(visible))

	// Confirm the two tools that the 33-tool bug was missing.
	mustHave := []string{"explode_xml_export", "generate_schema_map"}
	for _, name := range mustHave {
		if containsStr(visible, name) {
			fmt.Printf("        ✓ %s\n", name)
		} else {
			fmt.Printf("        ✗ %s  — MISSING. Your binary is stale; update it (Step 1).\n", name)
		}
	}

	if !tools.IsPluginConnected() {
		fmt.Println("        + 3 more (export_schema, read_layout, get_active_context)")
		fmt.Println("          appear once the FileMaker plugin connects (~35 total).")
	}
}

func containsStr(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// ── Version helpers ─────────────────────────────────────────────────────────

func fetchLatestVersion() (string, error) {
	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", latestAPIURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "kibuild-mcp-setup")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("GitHub returned status %d", resp.StatusCode)
	}
	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}
	if payload.TagName == "" {
		return "", fmt.Errorf("no tag_name in release")
	}
	return payload.TagName, nil
}

// isNewer reports whether latest > current using semver-ish comparison.
func isNewer(latest, current string) bool {
	if current == "dev" || current == "" {
		return true
	}
	lp := parseVer(latest)
	cp := parseVer(current)
	for i := 0; i < 3; i++ {
		if lp[i] != cp[i] {
			return lp[i] > cp[i]
		}
	}
	return false
}

func parseVer(v string) [3]int {
	var out [3]int
	v = strings.TrimPrefix(strings.TrimSpace(v), "v")
	// strip any pre-release/build suffix
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	parts := strings.Split(v, ".")
	for i := 0; i < 3 && i < len(parts); i++ {
		n, _ := strconv.Atoi(strings.TrimSpace(parts[i]))
		out[i] = n
	}
	return out
}

// ── Self-update ─────────────────────────────────────────────────────────────

// assetName returns the release asset filename for the current platform.
func assetName() (string, error) {
	switch runtime.GOOS {
	case "darwin":
		return "kibuild-mcp-darwin-" + runtime.GOARCH, nil
	case "linux":
		return "kibuild-mcp-linux-" + runtime.GOARCH, nil
	case "windows":
		return "kibuild-mcp-windows-amd64.exe", nil
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}

// selfUpdate downloads the latest binary for this platform and replaces the
// currently running executable. Works on all three OSes: on Unix the running
// file is replaced via rename; on Windows the running image is renamed aside
// first (allowed) and the new file moved into place.
func selfUpdate() error {
	asset, err := assetName()
	if err != nil {
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if resolved, rerr := filepath.EvalSymlinks(exe); rerr == nil {
		exe = resolved
	}

	dir := filepath.Dir(exe)
	tmp := filepath.Join(dir, ".kibuild-mcp.new")

	fmt.Printf("      Downloading %s ...\n", asset)
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(releaseBase + asset)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("cannot write to %s (try re-running the install script with elevated rights): %w", dir, err)
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		os.Remove(tmp)
		return err
	}
	out.Close()

	if runtime.GOOS == "windows" {
		old := exe + ".old"
		os.Remove(old) // ignore if absent
		if err := os.Rename(exe, old); err != nil {
			os.Remove(tmp)
			return fmt.Errorf("cannot move running binary aside: %w", err)
		}
		if err := os.Rename(tmp, exe); err != nil {
			// try to restore
			_ = os.Rename(old, exe)
			return fmt.Errorf("cannot place new binary: %w", err)
		}
		// best-effort cleanup; the .old file is released after this process exits
		return nil
	}

	// Unix: atomic replace of the running file.
	if err := os.Chmod(tmp, 0755); err != nil {
		os.Remove(tmp)
		return err
	}
	if err := os.Rename(tmp, exe); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("cannot replace %s (try the install script with sudo): %w", exe, err)
	}
	return nil
}

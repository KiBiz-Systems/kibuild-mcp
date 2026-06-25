package config

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Provider                 string `json:"provider,omitempty"`
	Model                    string `json:"model"`
	MaxTokens                int    `json:"maxTokens"`
	SQLTimeout               int    `json:"sqlTimeout"`
	MissionTimeout           int    `json:"mission_timeout"`
	QueryLimit               int    `json:"queryLimit"`
	PurgeToolOutput          bool   `json:"purge_tool_output"`
	PurgeThreshold           int    `json:"purge_threshold"`
	AutoSummarize            bool   `json:"auto_summarize"`
	SummarizeThreshold       int    `json:"summarize_threshold"`
	AntigravityCLIPermission string `json:"antigravity_cli_permission,omitempty"`
	ClaudeCLIPermission      string `json:"claude_cli_permission,omitempty"`
	CodexCLIPermission       string `json:"codex_cli_permission,omitempty"`
	AntigravityCLIPath       string `json:"antigravity_cli_path,omitempty"`
	ClaudeCLIPath            string `json:"claude_cli_path,omitempty"`
	CodexCLIPath             string `json:"codex_cli_path,omitempty"`
	AntigravityCLITimeout    int    `json:"antigravity_cli_timeout,omitempty"`
	ClaudeCLITimeout         int    `json:"claude_cli_timeout,omitempty"`
	CodexCLITimeout          int    `json:"codex_cli_timeout,omitempty"`
	OpenAIEnabled            bool   `json:"openai_enabled"`
	AnthropicEnabled         bool   `json:"anthropic_enabled"`
	GeminiEnabled            bool   `json:"gemini_enabled"`
	AntigravityCLIEnabled    bool   `json:"antigravity_cli_enabled"`
	ClaudeCLIEnabled         bool   `json:"claude_cli_enabled"`
	CodexCLIEnabled          bool   `json:"codex_cli_enabled"`
	AntigravityCLIMode       string `json:"antigravity_cli_mode,omitempty"`
	ClaudeCLIMode            string `json:"claude_cli_mode,omitempty"`
	CodexCLIMode             string `json:"codex_cli_mode,omitempty"`
	DisabledMCPTools         []string `json:"disabled_mcp_tools,omitempty"`
}

// SafeConfig is returned by the public /api/config endpoint.
// It intentionally omits all API key fields so credentials
// are never echoed back over HTTP.
type SafeConfig struct {
	Provider                 string `json:"provider,omitempty"`
	Model                    string `json:"model"`
	MaxTokens                int    `json:"maxTokens"`
	SQLTimeout               int    `json:"sqlTimeout"`
	MissionTimeout           int    `json:"mission_timeout"`
	QueryLimit               int    `json:"queryLimit"`
	PurgeToolOutput          bool   `json:"purge_tool_output"`
	PurgeThreshold           int    `json:"purge_threshold"`
	AutoSummarize            bool   `json:"auto_summarize"`
	SummarizeThreshold       int    `json:"summarize_threshold"`
	DetectedAntigravityCLI   bool   `json:"detected_antigravity_cli"`
	DetectedClaudeCLI        bool   `json:"detected_claude_cli"`
	DetectedCodexCLI         bool   `json:"detected_codex_cli"`
	AntigravityCLIPermission string `json:"antigravity_cli_permission,omitempty"`
	ClaudeCLIPermission      string `json:"claude_cli_permission,omitempty"`
	CodexCLIPermission       string `json:"codex_cli_permission,omitempty"`
	AntigravityCLIPath       string `json:"antigravity_cli_path,omitempty"`
	ClaudeCLIPath            string `json:"claude_cli_path,omitempty"`
	CodexCLIPath             string `json:"codex_cli_path,omitempty"`
	AntigravityCLITimeout    int    `json:"antigravity_cli_timeout,omitempty"`
	ClaudeCLITimeout         int    `json:"claude_cli_timeout,omitempty"`
	CodexCLITimeout          int    `json:"codex_cli_timeout,omitempty"`
	OpenAIEnabled            bool   `json:"openai_enabled"`
	AnthropicEnabled         bool   `json:"anthropic_enabled"`
	GeminiEnabled            bool   `json:"gemini_enabled"`
	AntigravityCLIEnabled    bool   `json:"antigravity_cli_enabled"`
	ClaudeCLIEnabled         bool   `json:"claude_cli_enabled"`
	CodexCLIEnabled          bool   `json:"codex_cli_enabled"`
	AntigravityCLIMode              string   `json:"antigravity_cli_mode,omitempty"`
	ClaudeCLIMode                   string   `json:"claude_cli_mode,omitempty"`
	CodexCLIMode                    string   `json:"codex_cli_mode,omitempty"`
	DisabledMCPTools                []string `json:"disabled_mcp_tools,omitempty"`
	DetectedAntigravityCLIPath      string   `json:"detected_antigravity_cli_path,omitempty"`
	DetectedClaudeCLIPath           string   `json:"detected_claude_cli_path,omitempty"`
	DetectedCodexCLIPath            string   `json:"detected_codex_cli_path,omitempty"`
}

type Manager struct {
	cfg  Config
	path string
}

func NewManager(configPath string) (*Manager, error) {
	m := &Manager{
		cfg: Config{
			Model:              "gemini-3.5-flash",
			MaxTokens:          4096,
			SQLTimeout:         5000,
			MissionTimeout:     900,
			QueryLimit:         500,
			PurgeToolOutput:    true,
			PurgeThreshold:     500,
			AutoSummarize:      true,
			SummarizeThreshold: 20000,
			OpenAIEnabled:          true,
			AnthropicEnabled:       true,
			GeminiEnabled:          true,
			AntigravityCLIEnabled:   true,
			ClaudeCLIEnabled:        true,
			CodexCLIEnabled:         true,
			AntigravityCLIMode:     "contract",
			ClaudeCLIMode:          "contract",
			CodexCLIMode:           "contract",
			DisabledMCPTools:       []string{},
		},
	}
	if configPath == "" {
		home, err := os.UserHomeDir()
		if err == nil {
			m.path = filepath.Join(home, ".fm_ai_bridge", "kibuild_settings.json")
		} else {
			m.path = filepath.Join("Config", "fm_ai_settings.json")
		}
	} else {
		m.path = configPath
	}

	if err := m.Load(m.path); err != nil {
		// If file doesn't exist, we continue using defaults + env vars.
	}
	m.loadEnvOverrides()
	return m, nil
}

func (m *Manager) ConfigPath() string {
	return m.path
}


func (m *Manager) Load(path string) error {
	m.path = path
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var temp map[string]interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Helper to extract string keys safely
	getString := func(key string) string {
		if val, ok := temp[key]; ok {
			if s, ok := val.(string); ok {
				return s
			}
		}
		return ""
	}

	getInt := func(key string, dflt int) int {
		if val, ok := temp[key]; ok {
			if n, ok := val.(float64); ok {
				return int(n)
			}
		}
		return dflt
	}

	getBool := func(key string, dflt bool) bool {
		if val, ok := temp[key]; ok {
			if b, ok := val.(bool); ok {
				return b
			}
		}
		return dflt
	}

	getSlice := func(key string) []string {
			if val, ok := temp[key]; ok {
				if list, ok := val.([]interface{}); ok {
					var res []string
					for _, item := range list {
						if s, ok := item.(string); ok {
							res = append(res, s)
						}
					}
					return res
				}
			}
			return []string{}
		}

	m.cfg.Provider = getString("provider")
	m.cfg.Model = getString("model")
	if m.cfg.Model == "" {
		m.cfg.Model = "gemini-3.5-flash"
	}
	m.cfg.MaxTokens = getInt("maxTokens", 4096)
	m.cfg.SQLTimeout = getInt("sqlTimeout", 5000)
	m.cfg.MissionTimeout = getInt("mission_timeout", 900)
	m.cfg.QueryLimit = getInt("queryLimit", 500)
	m.cfg.PurgeToolOutput = getBool("purge_tool_output", true)
	m.cfg.PurgeThreshold = getInt("purge_threshold", 500)
	m.cfg.AutoSummarize = getBool("auto_summarize", true)
	m.cfg.SummarizeThreshold = getInt("summarize_threshold", 20000)

	// CLI permissions
	m.cfg.AntigravityCLIPermission = getString("antigravity_cli_permission")
	m.cfg.ClaudeCLIPermission = getString("claude_cli_permission")
	m.cfg.CodexCLIPermission = getString("codex_cli_permission")

	// CLI paths and timeouts
	m.cfg.AntigravityCLIPath = getString("antigravity_cli_path")
	m.cfg.ClaudeCLIPath = getString("claude_cli_path")
	m.cfg.CodexCLIPath = getString("codex_cli_path")
	m.cfg.AntigravityCLITimeout = getInt("antigravity_cli_timeout", 0)
	m.cfg.ClaudeCLITimeout = getInt("claude_cli_timeout", 0)
	m.cfg.CodexCLITimeout = getInt("codex_cli_timeout", 0)

	m.cfg.OpenAIEnabled = getBool("openai_enabled", true)
	m.cfg.AnthropicEnabled = getBool("anthropic_enabled", true)
	m.cfg.GeminiEnabled = getBool("gemini_enabled", true)
	m.cfg.AntigravityCLIEnabled = getBool("antigravity_cli_enabled", true)
	m.cfg.ClaudeCLIEnabled = getBool("claude_cli_enabled", true)
	m.cfg.CodexCLIEnabled = getBool("codex_cli_enabled", true)

	m.cfg.AntigravityCLIMode = getString("antigravity_cli_mode")
	if m.cfg.AntigravityCLIMode == "" {
		m.cfg.AntigravityCLIMode = "contract"
	}
	m.cfg.ClaudeCLIMode = getString("claude_cli_mode")
	if m.cfg.ClaudeCLIMode == "" {
		m.cfg.ClaudeCLIMode = "contract"
	}
	m.cfg.CodexCLIMode = getString("codex_cli_mode")
	if m.cfg.CodexCLIMode == "" {
		m.cfg.CodexCLIMode = "contract"
	}
	m.cfg.DisabledMCPTools = getSlice("disabled_mcp_tools")

	return nil
}

func (m *Manager) loadEnvOverrides() {
	if val := os.Getenv("KIBUILD_MODEL"); val != "" {
		m.cfg.Model = val
	}
	if val := os.Getenv("KIBUILD_PROVIDER"); val != "" {
		m.cfg.Provider = val
	}
}

func (m *Manager) GetModel() string {
	return m.cfg.Model
}

func (m *Manager) GetMaxTokens() int {
	return m.cfg.MaxTokens
}

func (m *Manager) GetActiveProvider() string {
	if m.cfg.Provider != "" {
		prov := strings.ToLower(m.cfg.Provider)
		if prov == "antigravity_cli" || prov == "claude_cli" || prov == "codex_cli" {
			return m.cfg.Provider
		}
		model := strings.ToLower(m.cfg.Model)
		isGeminiModel := strings.HasPrefix(model, "gemini")
		isClaudeModel := strings.HasPrefix(model, "claude")
		isOpenAIModel := strings.HasPrefix(model, "gpt") || strings.HasPrefix(model, "o1") || strings.HasPrefix(model, "o3")

		if (prov == "gemini" && !isGeminiModel) ||
			(prov == "anthropic" && !isClaudeModel) ||
			(prov == "openai" && !isOpenAIModel) {
			// Mismatch! Fall through to infer from model
		} else {
			return m.cfg.Provider
		}
	}
	model := strings.ToLower(m.cfg.Model)
	switch {
	case model == "antigravity_cli" || model == "antigravity" || model == "agy":
		return "antigravity_cli"
	case model == "claude_cli":
		return "claude_cli"
	case model == "codex_cli" || model == "codex":
		return "codex_cli"
	case strings.HasPrefix(model, "claude"):
		return "anthropic"
	case strings.HasPrefix(model, "gemini"):
		return "gemini"
	default:
		return "openai"
	}
}

func getRealHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	if home == "" {
		return ""
	}
	if idx := strings.Index(home, "/Library/Containers/"); idx != -1 {
		home = home[:idx]
	}
	return home
}

func (m *Manager) GetConfig() Config {
	return m.cfg
}

// GetSafeConfig returns a redacted view of the config safe for HTTP exposure.
func (m *Manager) GetSafeConfig() SafeConfig {
	// --- Antigravity CLI detection ---
	var detectedAntiPath string
	if p, err := exec.LookPath("antigravity"); err == nil {
		detectedAntiPath = p
	} else if p, err := exec.LookPath("agy"); err == nil {
		detectedAntiPath = p
	} else {
		home := getRealHomeDir()
		if home != "" {
			for _, candidate := range []string{
				filepath.Join(home, ".local", "bin", "antigravity"),
				filepath.Join(home, ".local", "bin", "agy"),
				filepath.Join(home, ".gemini", "antigravity", "bin", "antigravity"),
				filepath.Join(home, ".gemini", "antigravity", "bin", "agy"),
			} {
				if _, err := os.Stat(candidate); err == nil {
					detectedAntiPath = candidate
					break
				}
			}
		}
	}

	// --- Claude CLI detection ---
	var detectedClaudePath string
	if p, err := exec.LookPath("claude"); err == nil {
		detectedClaudePath = p
	} else {
		home := getRealHomeDir()
		if home != "" {
			for _, candidate := range []string{
				filepath.Join(home, ".local", "bin", "claude"),
				filepath.Join(home, ".npm-global", "bin", "claude"),
				"/usr/local/bin/claude",
			} {
				if _, err := os.Stat(candidate); err == nil {
					detectedClaudePath = candidate
					break
				}
			}
		}
	}

	// --- Codex CLI detection ---
	var detectedCodexPath string
	if p, err := exec.LookPath("codex"); err == nil {
		detectedCodexPath = p
	} else {
		home := getRealHomeDir()
		if home != "" {
			for _, candidate := range []string{
				filepath.Join(home, ".gemini", "antigravity", "bin", "codex"),
				filepath.Join(home, ".gemini", "antigravity-ide", "bin", "codex"),
				filepath.Join(home, ".local", "bin", "codex"),
			} {
				if _, err := os.Stat(candidate); err == nil {
					detectedCodexPath = candidate
					break
				}
			}
		}
	}

	return SafeConfig{
		Provider:                 m.cfg.Provider,
		Model:                    m.cfg.Model,
		MaxTokens:                m.cfg.MaxTokens,
		SQLTimeout:               m.cfg.SQLTimeout,
		MissionTimeout:           m.cfg.MissionTimeout,
		QueryLimit:               m.cfg.QueryLimit,
		PurgeToolOutput:          m.cfg.PurgeToolOutput,
		PurgeThreshold:           m.cfg.PurgeThreshold,
		AutoSummarize:            m.cfg.AutoSummarize,
		SummarizeThreshold:       m.cfg.SummarizeThreshold,
		DetectedAntigravityCLI:         detectedAntiPath != "",
		DetectedClaudeCLI:              detectedClaudePath != "",
		DetectedCodexCLI:               detectedCodexPath != "",
		DetectedAntigravityCLIPath:     detectedAntiPath,
		DetectedClaudeCLIPath:          detectedClaudePath,
		DetectedCodexCLIPath:           detectedCodexPath,
		AntigravityCLIPermission: m.cfg.AntigravityCLIPermission,
		ClaudeCLIPermission:      m.cfg.ClaudeCLIPermission,
		CodexCLIPermission:       m.cfg.CodexCLIPermission,
		AntigravityCLIPath:       m.cfg.AntigravityCLIPath,
		ClaudeCLIPath:            m.cfg.ClaudeCLIPath,
		CodexCLIPath:             m.cfg.CodexCLIPath,
		AntigravityCLITimeout:    m.cfg.AntigravityCLITimeout,
		ClaudeCLITimeout:         m.cfg.ClaudeCLITimeout,
		CodexCLITimeout:          m.cfg.CodexCLITimeout,
		OpenAIEnabled:            m.cfg.OpenAIEnabled,
		AnthropicEnabled:         m.cfg.AnthropicEnabled,
		GeminiEnabled:            m.cfg.GeminiEnabled,
		AntigravityCLIEnabled:    m.cfg.AntigravityCLIEnabled,
		ClaudeCLIEnabled:         m.cfg.ClaudeCLIEnabled,
		CodexCLIEnabled:          m.cfg.CodexCLIEnabled,
		AntigravityCLIMode:       m.cfg.AntigravityCLIMode,
		ClaudeCLIMode:            m.cfg.ClaudeCLIMode,
		CodexCLIMode:             m.cfg.CodexCLIMode,
		DisabledMCPTools:         m.cfg.DisabledMCPTools,
	}
}

func (m *Manager) SaveConfig(cfg Config) error {
	// Clamp CLI mode values — only "mcp" and "contract" are valid.
	// Any unrecognised or empty string falls back to "contract".
	clampMode := func(mode string) string {
		if mode == "mcp" {
			return "mcp"
		}
		return "contract"
	}
	cfg.AntigravityCLIMode = clampMode(cfg.AntigravityCLIMode)
	cfg.ClaudeCLIMode = clampMode(cfg.ClaudeCLIMode)
	cfg.CodexCLIMode = clampMode(cfg.CodexCLIMode)

	m.cfg = cfg


	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	return os.WriteFile(m.path, data, 0600)
}


package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/priyabratasahoo21/kibuild-mcp/config"
	"github.com/priyabratasahoo21/kibuild-mcp/providers"
	"github.com/priyabratasahoo21/kibuild-mcp/skills"
	"github.com/priyabratasahoo21/kibuild-mcp/tools"
)

var cfgManager *config.Manager

func main() {
	// When run directly in a terminal (not spawned by an MCP client), print usage and exit.
	if fi, err := os.Stdin.Stat(); err == nil && (fi.Mode()&os.ModeCharDevice != 0) {
		fmt.Fprintln(os.Stderr, "kibuild-mcp — FileMaker MCP server")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "This binary is an MCP server. It is not run directly.")
		fmt.Fprintln(os.Stderr, "Register it in your AI tool's MCP config, then your AI tool will spawn it automatically.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example (Claude Code, ~/.claude.json):")
		fmt.Fprintln(os.Stderr, `  "mcpServers": { "kibuild": { "command": "/usr/local/bin/kibuild-mcp", "env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/project" } } }`)
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "See https://github.com/priyabratasahoo21/kibuild-mcp for setup instructions.")
		os.Exit(0)
	}

	var err error
	cfgManager, err = config.NewManager("")
	if err != nil {
		logMCPServer("Failed to initialize config: %v", err)
	}

	home, _ := os.UserHomeDir()
	skillsDir := filepath.Join(home, ".fm_ai_bridge", "skills")
	if errInit := skills.InitDefaultSkills(skillsDir); errInit != nil {
		logMCPServer("Warning: Failed to initialize default skills: %v", errInit)
	}
	tools.SetSkillsDir(skillsDir)

	runMCPLoop()
}

// MCP JSON-RPC types

type mcpRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type mcpResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}

func runMCPLoop() {
	logMCPServer("kibuild-mcp started, pid=%d", os.Getpid())
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			logMCPServer("stdin read error: %v, exiting", err)
			break
		}

		lineStr := strings.TrimSpace(string(line))
		logMCPServer("Incoming: %s", lineStr)

		var req mcpRequest
		if err := json.Unmarshal(line, &req); err != nil {
			logMCPServer("JSON parse error: %v", err)
			sendMCPError(nil, -32700, "Parse error")
			continue
		}

		handleMCPRequest(&req)
	}
}

func handleMCPRequest(req *mcpRequest) {
	logMCPServer("Handling: method=%s, id=%v", req.Method, req.ID)

	if req.ID == nil || strings.HasPrefix(req.Method, "notifications/") {
		return
	}

	switch req.Method {
	case "initialize":
		sendMCPResponse(req.ID, map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]string{
				"name":    "kibuild-mcp",
				"version": "1.0.0",
			},
		})

	case "tools/list":
		type mcpToolDef struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			InputSchema interface{} `json:"inputSchema"`
		}

		rawTools := tools.GetToolsSchema()
		var mcpTools []mcpToolDef
		for _, t := range rawTools {
			if !tools.IsSafeMCPTool(t.Name) || isToolDisabled(t.Name) {
				continue
			}
			mcpTools = append(mcpTools, mcpToolDef{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: t.Parameters,
			})
		}
		sendMCPResponse(req.ID, map[string]interface{}{"tools": mcpTools})

	case "tools/call":
		var params struct {
			Name      string          `json:"name"`
			Arguments json.RawMessage `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			sendMCPError(req.ID, -32602, "Invalid params")
			return
		}

		logMCPServer("tools/call: %s args=%s", params.Name, string(params.Arguments))

		if !tools.IsSafeMCPTool(params.Name) || isToolDisabled(params.Name) {
			sendMCPError(req.ID, -32601, fmt.Sprintf("Tool %s is blocked or disabled", params.Name))
			return
		}

		ctx := context.Background()
		activeProj := os.Getenv("KIBUILD_ACTIVE_PROJECT")
		if activeProj == "" {
			if home, err := os.UserHomeDir(); err == nil {
				if data, err := os.ReadFile(filepath.Join(home, ".fm_ai_bridge", "active_project.txt")); err == nil {
					activeProj = strings.TrimSpace(string(data))
				}
			}
		}
		if activeProj != "" {
			ctx = context.WithValue(ctx, providers.ProjectPathKey, activeProj)
		}

		out, err := tools.ExecuteTool(ctx, params.Name, string(params.Arguments))
		if err != nil {
			logMCPServer("tool %s error: %v", params.Name, err)
			sendMCPResponse(req.ID, map[string]interface{}{
				"isError": true,
				"content": []map[string]interface{}{{"type": "text", "text": fmt.Sprintf("Error: %v", err)}},
			})
			return
		}

		logMCPServer("tool %s ok (%d bytes)", params.Name, len(out))
		sendMCPResponse(req.ID, map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": out}},
		})

	default:
		logMCPServer("Unknown method: %s", req.Method)
		sendMCPError(req.ID, -32601, "Method not found")
	}
}

func isToolDisabled(name string) bool {
	if cfgManager != nil {
		_ = cfgManager.Load(cfgManager.ConfigPath())
		for _, t := range cfgManager.GetConfig().DisabledMCPTools {
			if t == name {
				return true
			}
		}
	}
	return false
}

func sendMCPResponse(id interface{}, result interface{}) {
	resp := mcpResponse{JSONRPC: "2.0", Result: result, ID: id}
	b, _ := json.Marshal(resp)
	logMCPServer("Outgoing: %s", string(b))
	fmt.Printf("%s\n", string(b))
}

func sendMCPError(id interface{}, code int, message string) {
	resp := mcpResponse{
		JSONRPC: "2.0",
		Error:   map[string]interface{}{"code": code, "message": message},
		ID:      id,
	}
	b, _ := json.Marshal(resp)
	logMCPServer("Error: %s", string(b))
	fmt.Printf("%s\n", string(b))
}

func logMCPServer(format string, args ...interface{}) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	logPath := filepath.Join(home, ".fm_ai_bridge", "mcp_server.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer f.Close()
	msg := fmt.Sprintf(format, args...)
	f.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format("2006-01-02 15:04:05.000"), msg))
}

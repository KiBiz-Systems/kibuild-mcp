# KiBuild MCP — Claude Code Project Guide

## Magic phrase

If the user types **"Help me set up KiBuild MCP"** or **"Set up KiBuild MCP"**, run the `/setup-kibuild` wizard automatically: detect the OS, install the binary, configure the MCP server in `~/.claude.json`, verify the tool count, and offer to test with an XML export. Follow the steps in `.claude/commands/setup-kibuild.md` exactly.

## Project initialization command

Run `/init-kibuild-project` inside any FileMaker project directory to generate three AI guide files:
- **`CLAUDE.md`** — instructs Claude Code to prefer KiBuild MCP tools over grep/bash
- **`AGENTS.md`** — same, for OpenAI Codex CLI and ChatGPT Agents
- **`GEMINI.md`** — same, for Google Gemini CLI

Each file contains explicit tool-substitution rules (e.g. "use `find_script` instead of grep") so the AI picks the right tool without being told each time. Command lives at `.claude/commands/init-kibuild-project.md`.

## Native setup command

The binary has a built-in interactive wizard — `kibuild-mcp --setup` — that checks the version (self-updates from GitHub if behind), prompts for the project path, safely merges the MCP entry into `~/.claude.json`, and verifies the tool registry. The install scripts hand off to it after downloading the binary. All config/verify logic lives in [setup.go](setup.go), so it behaves identically on macOS, Linux, and Windows. `kibuild-mcp --version` prints the version (injected at build time via `-ldflags "-X main.Version=..."`; defaults to `dev`).

## What this project is

KiBuild MCP is a FileMaker-aware MCP server. It provides AI coding tools with structured access to FileMaker schemas, scripts, layouts, relationships, and XML exports.

## Repository layout

| Path | Purpose |
|---|---|
| `main.go` | MCP JSON-RPC server loop and tool dispatch |
| `tools/registry.go` | Tool definitions + `safeMCPTools` allowlist |
| `tools/executor.go` | Tool implementations |
| `tools/discovery.go` | Schema discovery helpers |
| `exploder/` | XML export exploder (produces `Schema/<db>/` tree) |
| `config/` | Config manager, sandbox, crypto |
| `outbox/` | Versioned artifact output manager |
| `skills/` | Embedded specialist skill files |
| `install.sh` | macOS/Linux installer |
| `install.ps1` | Windows PowerShell installer |
| `.claude/commands/setup-kibuild.md` | `/setup-kibuild` slash command |

## Tool count reference

KiBuild MCP exposes **32 tools**. If the count is under 28, the binary is likely outdated. Ask the user to reinstall and share the last 20 lines of `~/.fm_ai_bridge/mcp_server.log`.

## Key tools (v0.2.0+)

- `explode_xml_export` — converts a FileMaker XML export into the per-object `Schema/` tree
- `generate_schema_map` — indexes the Schema tree into `workspace_map.md` for fast RAG
- `find_script` / `find_table` / `find_layout` — schema lookup
- `search_index` — token-efficient keyword search over `workspace_map.md`
- `write_outbox_artifact` — writes versioned generated output to the project outbox

## Debugging missing tools

1. Run `kibuild-mcp --setup` — step 4 prints the exact tool count and confirms whether `explode_xml_export` / `generate_schema_map` are present. This is the fastest diagnosis.
2. Check binary version: `kibuild-mcp --version` (must be v0.2.0+).
3. Check server log: `tail -20 ~/.fm_ai_bridge/mcp_server.log`
4. Confirm `explode_xml_export` and `generate_schema_map` are in `safeMCPTools` in `tools/registry.go` (they are in v0.2.0+).
5. If outdated: accept the self-update in `kibuild-mcp --setup`, or reinstall with `curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh`. Then **fully quit and restart Claude Code** — the MCP client caches the old process.

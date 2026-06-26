# KiBuild MCP

A FileMaker-aware MCP server that gives AI coding tools deep, structured access to your FileMaker schema — scripts, layouts, tables, relationships, and dependencies — without grep, without reading raw XML, and without a FileMaker license at runtime.

Works with Claude Code, Cursor, Windsurf, Codex, Antigravity (Google), and any MCP-compatible AI tool.

---

## Why KiBuild instead of grep?

When an AI uses `grep` or reads raw XML to answer a FileMaker question, it is flying blind. Here is what that looks like in practice:

| Task | Without KiBuild | With KiBuild |
|---|---|---|
| Find where a script is called | `grep -r "Send Email"` across hundreds of XML files — noisy, slow, misses numeric IDs | `find_script_references_in_scripts` returns exact call sites with file paths and line numbers |
| Find where a field is used | Manual regex over calculation XML — misses layouts and relationships | `find_field_references_in_calculations / _in_layouts / _in_relationships` in one call |
| Understand a script | Read hundreds of lines of raw XML step-by-step | `find_script` returns a clean step list with names resolved |
| Trace a relationship | Open the relationship graph XML and parse join predicates manually | `inspect_relationships` returns structured JSON immediately |
| Search the whole schema | `grep` tokenizes everything — expensive, context-filling | `search_index` searches `workspace_map.md` — token-efficient, instant |

**Who should use this:**
- FileMaker developers using AI tools to navigate, understand, or generate schema changes
- Teams doing AI-assisted code review, impact analysis, or script migrations
- Anyone who has asked an AI "where is this field used?" and gotten a slow, incomplete answer

---

## How it works

```
Your FileMaker file
      │
      │  Tools → Save a Copy as XML…  (or via script step)
      ▼
XML export  ──→  explode_xml_export  ──→  Schema/<db>/
                                                    scripts/
                                                    layouts/
                                                    tables/
                                                    relationships/
                                                    (one file per object)
                                                         │
                                                         ▼
                                               generate_schema_map
                                                         │
                                                         ▼
                                               workspace_map.md
                                          (compact index, fast AI search)
                                                         │
                                                         ▼
                                           32 KiBuild tools live and ready
```

Export your schema from FileMaker using **Tools → Save a Copy as XML…** or a script step — see Step 2 for full instructions. The exploder auto-detects the export format and processes it into one file per object.

---

## Get started

### Step 1 — Install the binary

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex
```

> **No `curl`?** The script falls back to `wget` automatically.
>
> **Windows execution policy error?** Run this first:
> ```powershell
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```

The installer downloads the binary, runs `kibuild-mcp --setup` to configure your AI tool's MCP config, and installs AI guide files globally.

---

### Step 2 — Export your FileMaker schema

You need to export your schema from FileMaker once. After that, the exploder keeps it as a searchable file tree.

#### Method 1 — Via menu (interactive)

**FileMaker 2025:**
1. Open your FileMaker file
2. Go to **Tools → Save a Copy as XML…**
3. Choose a destination folder and filename
4. Optionally check **Include details for analysis tools** to embed the DDR_INFO block (larger file, useful for deeper analysis)
5. Click **Save** — outputs a single `.xml` file covering the entire solution

**FileMaker 2026:**
1. Open the file(s) you want to export
2. Go to **Tools → Save a Copy as XML…**
3. A dialog shows a checklist of the ~20 available catalogs (Scripts, Layouts, Base Tables, Custom Functions, etc.) — check only what you need
4. Optionally split binary layout data into a separate catalog file
5. Optionally export multiple open files in a single operation
6. Click **Save** — outputs one or more targeted `.xml` files

#### Method 2 — Via script step (automated, both versions)

```
Save a Copy as XML [ Window Name: "YourFileName" ; Destination File: "/path/to/output.xml" ]
```

- **FileMaker 2025** — supports `Window Name`, `Destination File`, and `Include details for analysis tools`
- **FileMaker 2026** — same base parameters, extended with additional options for catalog selection and splitting, enabling fully scripted targeted exports

Then ask your AI tool to process it:
```
Explode the export at /path/to/MyDatabase.xml into my project, then build the schema map.
```

The AI calls `explode_xml_export` (splits into `Schema/<db>/`) and `generate_schema_map` (builds `workspace_map.md`). After this step, all navigation and reference tools are live.

> **What gets exploded:** scripts (+ plain-text sanitized copy), tables with fields, layouts, relationships, table occurrences, value lists, custom functions, custom menus, accounts, privilege sets, themes — one file per object, ready for Git diffing.

---

### Step 3 — Initialize your project and start working

Open your FileMaker project folder in your AI tool and run:

**Claude Code:**
```
/init-kibuild-project
```

This writes `CLAUDE.md`, `AGENTS.md`, `GEMINI.md`, and `.cursor/rules/kibuild.mdc` to the project root — telling whichever AI tool opens the folder to prefer KiBuild tools over grep. Commit these files so your whole team gets the same behavior automatically.

Then ask natural questions:
```
Find the script "Create Invoice" and show me what it does.
Which scripts call "Send Email Notification"?
Where is the Status field used across scripts, layouts, and calculations?
List all layouts that reference the Invoices table occurrence.
Show me the relationships for the Contacts table.
```

---

### One phrase, zero commands (Claude Code only)

Open Claude Code in this repo's folder and type:

```
Help me set up KiBuild MCP
```

Claude reads the project guide and runs the full setup wizard automatically. No terminal required.

---

## Configure for your AI tool

The setup wizard handles Claude Code automatically. For other tools, add the MCP server entry manually. Replace `/path/to/your/project` with the absolute path to your FileMaker project folder.

### Claude Code

**macOS / Linux** — `~/.claude.json`
```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

**Windows** — `C:\Users\<YourName>\.claude.json`
```json
{
  "mcpServers": {
    "kibuild": {
      "command": "C:/Users/<YourName>/AppData/Local/Programs/kibuild-mcp/kibuild-mcp.exe",
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "C:/Users/<YourName>/Documents/MyFileMakerProject"
      }
    }
  }
}
```

After editing, restart Claude Code or run `/mcp` to reload.

---

### Cursor

Config file: `~/.cursor/mcp.json`
```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project" }
    }
  }
}
```

---

### Windsurf

Config file: `~/.codeium/windsurf/mcp_config.json`
```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project" }
    }
  }
}
```

---

### OpenAI Codex CLI

Config file: `~/.codex/config.toml`
```toml
[mcp_servers.kibuild]
command = "/usr/local/bin/kibuild-mcp"

[mcp_servers.kibuild.env]
KIBUILD_ACTIVE_PROJECT = "/path/to/your/project"
```

---

### Google Antigravity (Agy)

Config file: `~/.gemini/config/mcp_config.json`
```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project" }
    }
  }
}
```

Create the folder if needed: `mkdir -p ~/.gemini/config`

---

### VS Code (with MCP extension)

User `settings.json` — open via `Ctrl+Shift+P` → "Open User Settings JSON"
```json
{
  "mcp.servers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project" }
    }
  }
}
```

---

## Useful commands

| Command | Where | What it does |
|---|---|---|
| `kibuild-mcp --setup` | Terminal | Full wizard: version check, config, tool verification, installs AI guides |
| `kibuild-mcp --version` | Terminal | Print the installed version |
| `/setup-kibuild` | Claude Code | Same wizard driven by Claude Code |
| `/init-kibuild-project` | Claude Code | Write AI guide files to the current project folder |
| `/mcp` | Claude Code | List connected MCP servers and their tool count |

---

## Tool reference

### Schema navigation

| Tool | What it does |
|---|---|
| `find_script` | Find a script by name — returns sanitized step list, file paths, and sibling scripts |
| `find_layout` | Find a layout — returns bound table occurrence, referenced scripts, and XML path |
| `find_table` | Find a base table — returns all fields with types and XML path |
| `inspect_relationships` | Return all relationship predicates for a database or table occurrence |
| `search_index` | Keyword search over `workspace_map.md` — token-efficient, returns only matching lines |
| `generate_schema_map` | Build or refresh `workspace_map.md` — compact index of all schema objects |

### Impact analysis

| Tool | What it finds |
|---|---|
| `find_script_references_in_scripts` | Where a script is called via Perform Script |
| `find_script_references_in_layouts` | Layouts that trigger a script via buttons or script triggers |
| `find_script_references_to_layouts` | Go to Layout steps inside a script |
| `find_script_references_to_valuelists` | Value list references inside a script |
| `find_field_references_in_scripts` | Scripts that read or write a field |
| `find_field_references_in_layouts` | Layouts that display a field |
| `find_field_references_in_calculations` | Calc fields, auto-enter calcs, and validation rules referencing a field |
| `find_field_references_in_relationships` | Relationship join predicates that use a field |
| `find_variable_references_in_scripts` | Scripts that set or read a `$variable` |
| `find_layout_references_to_scripts` | Scripts triggered by a layout's buttons or script triggers |
| `find_layout_references_to_valuelists` | Value lists used by field controls on a layout |
| `find_layout_references_to_tables` | Table occurrences referenced by fields on a layout |
| `find_layout_references_in_calculations` | Calculations that reference a layout name |
| `find_valuelist_references_in_calculations` | Calculations that reference a value list |
| `find_to_references` | Every layout, script, and relationship referencing a table occurrence |
| `find_relationship_predicates` | Full join predicate details for a table occurrence |

### XML and generation

| Tool | What it does |
|---|---|
| `explode_xml_export` | Split a FileMaker XML export into the per-object schema tree. Auto-detects format. |
| `xml_extract_steps` | List all script steps from an FMXML snippet or file |
| `xml_lookup_name` | Resolve a numeric script ID to its name |
| `xml_trace_dependencies` | Extract all referenced TOs, scripts, layouts, and fields from XML |
| `xml_match_revision` | Read FileMaker version and revision metadata from an XML header |
| `validate_fmxmlsnippet` | Run 7-rule structural validation on a generated FMXML snippet |
| `validate_webviewer_html` | Check WebViewer HTML for remote dependencies and risky APIs |
| `write_outbox_artifact` | Save generated output to the project outbox as a versioned artifact |
| `read_xml_guide` | Load the FileMaker XML snippet reference guide — schemas and templates for script step generation |

### Specialist skills

| Tool | What it does |
|---|---|
| `load_skill` | Load a specialist skill into AI context. Available: `pro_scriptwriter`, `script_analysis`, `fm_xml_serializer`, `script_debug` |

---

## Troubleshooting

**Fewer than 30 tools visible?**

Your binary is outdated (pre-v0.2.0). Reinstall:
```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh

# Windows
irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex
```

Then fully quit and restart your AI tool — the MCP client caches the old process.

**macOS binary blocked by Gatekeeper?**
```bash
xattr -d com.apple.quarantine /usr/local/bin/kibuild-mcp
```

Then allow it in **System Settings → Privacy & Security** if prompted.

**Server not starting?**

Check the log:
```bash
# macOS / Linux
tail -20 ~/.fm_ai_bridge/mcp_server.log

# Windows
Get-Content "$env:USERPROFILE\.fm_ai_bridge\mcp_server.log" -Tail 20
```

Common causes: binary path mismatch in config, Gatekeeper blocking on macOS, PATH issues on Windows.

---

## Tool count reference

KiBuild MCP exposes **32 tools**. If you see fewer than 28, the binary is likely outdated — run `kibuild-mcp --setup` to self-update, then fully restart Claude Code.

---

## Installing from source

**go install** (requires Go 1.21+):
```bash
go install github.com/priyabratasahoo21/kibuild-mcp@latest
kibuild-mcp --setup
```

**Build from source:**
```bash
git clone https://github.com/priyabratasahoo21/kibuild-mcp.git
cd kibuild-mcp
go build -ldflags="-s -w -X main.Version=v0.2.0" -o kibuild-mcp .
sudo mv kibuild-mcp /usr/local/bin/
kibuild-mcp --setup
```

**Manual download** — [Releases page](https://github.com/priyabratasahoo21/kibuild-mcp/releases)

| Platform | File |
|---|---|
| macOS (Apple Silicon) | `kibuild-mcp-darwin-arm64` |
| macOS (Intel) | `kibuild-mcp-darwin-amd64` |
| Linux (x86_64) | `kibuild-mcp-linux-amd64` |
| Linux (ARM64) | `kibuild-mcp-linux-arm64` |
| Windows | `kibuild-mcp-windows-amd64.exe` |

---

## Architecture

```
AI tool (Claude Code, Cursor, Windsurf, VS Code, Codex, Antigravity)
  │
  │  spawns subprocess on MCP connect
  ▼
kibuild-mcp binary
  │  MCP JSON-RPC over stdin/stdout
  │
  └── 32 tools reading Schema/ XML files on disk
        no FileMaker license required at runtime

Schema location:  ~/your-project/files/Schema/<DBName>/
Project index:    ~/your-project/workspace_map.md
Log:              ~/.fm_ai_bridge/mcp_server.log
```

---

## Contributing

Pull requests welcome. Open an issue first for significant changes.

## License

MIT

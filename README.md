# KiBuild MCP

A self-contained MCP server that gives any AI coding tool deep, read-level access to a Claris FileMaker schema — no FileMaker license required at runtime.

Register one binary. Your AI tool (Claude Code, Cursor, Windsurf, Codex, Antigravity) gains FileMaker-aware tools: script navigation, full dependency graph, XML validation, and specialist skills.

---

## Get started in 60 seconds

Run the installer for your platform. It downloads the binary and launches a setup wizard that configures everything automatically.

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex
```

The wizard will:
1. Check for the latest version and self-update if needed
2. Ask for your FileMaker project folder
3. Write the MCP config for your AI tool
4. Verify all tools are accessible and print the tool count

Then **restart your AI tool**. That's it.

> **No `curl`?** The script falls back to `wget` automatically.
>
> **Windows execution policy error?** Run this once first:
> ```powershell
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```

---

### One phrase, zero commands

**Claude Code** — open Claude Code with this repo's folder and type:

```
Help me set up KiBuild MCP
```

Claude reads the project guide and runs the full wizard automatically. No terminal needed.

**Any other AI tool** — paste this into the chat:

```
Set up KiBuild MCP from: https://github.com/priyabratasahoo21/kibuild-mcp
```

The AI reads the README, walks you through the installer, and configures the MCP entry for your tool.

---

## What it does

### Schema navigation

Find any script, layout, or table by name with fuzzy matching. Returns the sanitized step list, sibling scripts, and the raw XML path — everything the AI needs to reason about the schema without parsing XML itself.

### Full dependency graph

Trace anything to anything: which layouts trigger a script, which scripts navigate to a layout, where a field is used in calculations or join predicates, which value lists appear in layout controls. 16 reference tools walk the exploded schema XML and return structured JSON with file paths and line-level snippets.

### XML analysis and generation

Extract and list script steps from XML, validate generated FMXML snippets against 7 structural rules before they reach FileMaker, validate WebViewer HTML for remote dependencies and risky APIs, and write versioned artifacts to the project outbox for review.

### Specialist skills

Load curated FileMaker skill prompts (`pro_scriptwriter`, `script_analysis`, `fm_xml_serializer`, `script_debug`) directly into AI context to inject domain-specific guidance for writing, analyzing, or debugging scripts.

---

## Configure for your AI tool

After the binary is installed, add the MCP server entry to your tool's config file. Replace `/path/to/your/project` with the absolute path to your FileMaker project folder (the folder that contains `files/Schema/`).

> **Tip:** The `kibuild-mcp --setup` wizard does this for you automatically for Claude Code. For other tools, use the snippets below.

---

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
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
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
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

---

### OpenAI Codex CLI

Config file: `~/.codex/config.toml` (global) or `.codex/config.toml` (project-scoped)

```toml
[mcp_servers.kibuild]
command = "/usr/local/bin/kibuild-mcp"

[mcp_servers.kibuild.env]
KIBUILD_ACTIVE_PROJECT = "/path/to/your/project"
```

Or add via CLI:
```bash
codex mcp add kibuild -- /usr/local/bin/kibuild-mcp
```
Then open the config and add `KIBUILD_ACTIVE_PROJECT` to the env block manually.

---

### Google Antigravity (Agy)

Config file: `~/.gemini/config/mcp_config.json`

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

Create the folder if needed: `mkdir -p ~/.gemini/config`

---

### VS Code (with MCP extension)

User `settings.json` — open via `Ctrl+Shift+P` → "Open User Settings JSON"

```json
{
  "mcp.servers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

---

## Get your FileMaker schema

The server indexes an **exploded** schema folder — one XML file per object, grouped into `scripts/`, `layouts/`, `tables/`, and `relationships/`:

```
your-project/
└── files/
    └── Schema/
        └── YourDatabase/
            ├── scripts/
            ├── scripts_sanitized/
            ├── layouts/
            ├── tables/
            └── relationships/
```

### Schema sources

| Source | Output | Ready to index? |
|---|---|---|
| KiBuild plugin **Export Schema** | Exploded tree, one file per object | ✅ Directly |
| **Save a Copy as XML** — single file | One `FMSaveAsXML` document | ✅ Via `explode_xml_export` |
| **Save a Copy as XML** — per-catalog option | One file per catalog | ✅ Via `explode_xml_export` |
| Built-in **DDR** export | One `FMPReport` document | Not yet supported |

### Using `explode_xml_export`

If you exported using FileMaker's **Save a Copy as XML** (either as a single file or with the per-catalog split), use the built-in `explode_xml_export` tool to convert it into the per-object layout. Ask your AI tool:

```
Explode the Save-as-XML export at /path/to/Contacts.xml into my project, then build the schema map.
```

It auto-detects the format and writes one file per object under `Schema/<database>/`, then you run `generate_schema_map` to index it.

> **What gets exploded:** scripts (+ sanitized `.txt`), tables (fields joined in), layouts, relationships, table occurrences, value lists, custom functions, custom menus, accounts, privilege sets, extended privileges, themes, and more — one file per object, ready for Git diffing.

---

## Build the workspace index

Once the server is configured, ask your AI tool to index your schema:

```
Call generate_schema_map for my project at /path/to/your/project
```

This writes `workspace_map.md` to your project root. After that, all navigation and reference tools are live.

---

## Usage examples

Once set up, ask your AI tool natural questions:

```
Find the script "Create Invoice" and show me what it does.
```
```
Which scripts call "Send Email Notification"?
```
```
Where is the Status field used across scripts, layouts, and calculations?
```
```
List all layouts that reference the Invoices table occurrence.
```
```
Show me the relationships for the Contacts table occurrence.
```
```
Validate this FMXML snippet before I import it.
```

---

## Reference

### Useful commands

| Command | Where | What it does |
|---|---|---|
| `kibuild-mcp --setup` | Terminal | Full wizard: version check, config, tool verification |
| `kibuild-mcp --version` | Terminal | Print the installed version |
| `/setup-kibuild` | Claude Code | Same wizard driven by Claude Code, with extra diagnosis |
| `/mcp` | Claude Code | List connected MCP servers and their tools |

### Tool count

| State | Expected tools |
|---|---|
| No plugin connected | ~32 |
| Plugin connected | ~35 |

If you see fewer than 30 tools, the binary is likely outdated — run `kibuild-mcp --setup` or reinstall.

---

## Tool reference

### Schema navigation

| Tool | Description |
|---|---|
| `find_script` | Find a script by name. Returns sanitized step list, `txt_path`, `xml_path`, and sibling scripts. |
| `find_layout` | Find a layout by name. Returns bound table occurrence, referenced scripts and layouts, and the XML path. |
| `find_table` | Find a base table by name. Returns all fields with types and the XML path. |
| `inspect_relationships` | Return all relationship predicates for a database or table occurrence. |
| `search_index` | Keyword search over `workspace_map.md`. Token-efficient — returns only matching lines. |
| `generate_schema_map` | Build or refresh `workspace_map.md` — a compact index of all tables, layouts, scripts, and table occurrences. |

### Impact analysis — reference finding

| Tool | What it finds |
|---|---|
| `find_layout_references_to_scripts` | Scripts triggered by buttons or script triggers on the given layouts |
| `find_layout_references_to_valuelists` | Value lists used by field controls on the given layouts |
| `find_layout_references_to_tables` | Table occurrences referenced by fields on the given layouts |
| `find_script_references_in_scripts` | Locations where the given scripts are called via Perform Script |
| `find_script_references_in_layouts` | Layouts that trigger the given scripts via buttons or script triggers |
| `find_script_references_to_layouts` | Go to Layout steps inside the given scripts |
| `find_script_references_to_valuelists` | Value list references inside the given scripts |
| `find_field_references_in_scripts` | Scripts that read or write the given fields |
| `find_field_references_in_layouts` | Layouts that display the given fields |
| `find_field_references_in_calculations` | Calc fields, auto-enter calcs, and validation rules that reference the given fields |
| `find_field_references_in_relationships` | Relationship join predicates that use the given fields |
| `find_variable_references_in_scripts` | Scripts that set or read the given `$variable` names |
| `find_valuelist_references_in_calculations` | Calculations that reference the given value lists |
| `find_layout_references_in_calculations` | Calculations that reference the given layout names |
| `find_to_references` | Every layout, script, and relationship that references the given table occurrences |
| `find_relationship_predicates` | Full join predicate details (left/right TO, field, operator) for the given table occurrences |

### XML analysis and generation

| Tool | Description |
|---|---|
| `explode_xml_export` | Explode a FileMaker Save-as-XML export into the per-object schema layout. Auto-detects single-file or per-catalog format. |
| `xml_extract_steps` | List all script steps from a raw FMXML snippet or file content. |
| `xml_lookup_name` | Resolve a numeric script ID to its name from an XML document. |
| `xml_trace_dependencies` | Extract all referenced table occurrences, scripts, layouts, and fields from XML content. |
| `xml_match_revision` | Read the FileMaker version and revision metadata from an XML header. |
| `validate_fmxmlsnippet` | Run 7-rule structural validation on a generated FMXML snippet and return a pass/fail report. |
| `validate_webviewer_html` | Check generated WebViewer HTML for remote dependencies, risky JavaScript APIs, and FileMaker bridge usage. |
| `write_outbox_artifact` | Save a generated script, layout, or document to the project outbox as a versioned artifact. |

### Specialist skills

| Tool | Description |
|---|---|
| `load_skill` | Load a specialist skill into AI context. Available: `pro_scriptwriter`, `script_analysis`, `fm_xml_serializer`, `script_debug`. |

---

## Troubleshooting

**Seeing only 33 tools (missing `explode_xml_export` / `generate_schema_map`)?**

Your binary is pre-v0.2.0. Reinstall:

```bash
# macOS/Linux
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh

# Windows
irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex
```

Then fully quit and restart your AI tool (the MCP client caches the old process). Confirm with `kibuild-mcp --version`.

**macOS binary blocked by Gatekeeper?**

```bash
xattr -d com.apple.quarantine /usr/local/bin/kibuild-mcp
```

Then allow it in **System Settings → Privacy & Security** if prompted.

---

## Logging

The server logs all MCP traffic to `~/.fm_ai_bridge/mcp_server.log`:

```bash
# macOS / Linux
tail -f ~/.fm_ai_bridge/mcp_server.log

# Windows (PowerShell)
Get-Content -Wait "$env:USERPROFILE\.fm_ai_bridge\mcp_server.log"
```

---

## Installing from source

**Option A — `go install`** (requires Go 1.21+)

```bash
go install github.com/priyabratasahoo21/kibuild-mcp@latest
kibuild-mcp --setup
```

> `go install` reports version as `dev` — the `--setup` wizard will offer to fetch the latest release binary to get a version-stamped build.

**Option B — Build from source** (requires Go 1.21+)

```bash
git clone https://github.com/priyabratasahoo21/kibuild-mcp.git
cd kibuild-mcp
go build -ldflags="-s -w -X main.Version=v0.2.0" -o kibuild-mcp .
mv kibuild-mcp /usr/local/bin/
kibuild-mcp --setup
```

**Option C — Manual download**

Go to the [Releases page](https://github.com/priyabratasahoo21/kibuild-mcp/releases) and download the binary for your platform:

| Platform | File |
|---|---|
| macOS (Apple Silicon) | `kibuild-mcp-darwin-arm64` |
| macOS (Intel) | `kibuild-mcp-darwin-amd64` |
| Linux (x86_64) | `kibuild-mcp-linux-amd64` |
| Linux (ARM64) | `kibuild-mcp-linux-arm64` |
| Windows | `kibuild-mcp-windows-amd64.exe` |

```bash
chmod +x kibuild-mcp-*
sudo mv kibuild-mcp-* /usr/local/bin/kibuild-mcp
kibuild-mcp --setup
```

---

## Architecture

```
AI tool (Claude Code, Cursor, Windsurf, VS Code, Codex, Antigravity)
  │
  │  spawns subprocess on MCP connect
  ▼
kibuild-mcp  ← this binary
  │  MCP JSON-RPC over stdin/stdout (protocol 2024-11-05)
  │
  └── analysis tools  ← read Schema/ XML files on disk
        works without FileMaker running

Reads from:
  ~/your-project/files/Schema/<DBName>/   ← exported schema (XML files)
  ~/.fm_ai_bridge/active_project.txt      ← current project pointer
  ~/your-project/workspace_map.md         ← built by generate_schema_map
```

---

## Contributing

Pull requests are welcome. Please open an issue first for significant changes.

---

## License

MIT

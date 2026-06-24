# KiBuild MCP

A self-contained MCP server that gives any AI coding tool deep, read-level access to a Claris FileMaker schema — no FileMaker license required at runtime for schema analysis.

Register one binary in your MCP config. Your AI tool (Claude Code, Cursor, Windsurf, VS Code) immediately gains FileMaker-aware tools: script navigation, impact analysis across the full dependency graph, XML generation and validation, and specialist skills.

---

## What it does

### Schema navigation
Find any script, layout, or table by name with fuzzy matching. Returns the sanitized step list, sibling scripts, and the raw XML path — everything the AI needs to reason about the schema without reading raw XML itself.

### Full dependency graph (16 reference tools)
Trace anything to anything: which layouts trigger a script, which scripts navigate to a layout, where a field is used in calculations or join predicates, which value lists appear in layout controls. Every reference tool walks the exploded schema XML and returns structured JSON with file paths and line-level snippets.

### XML analysis and generation
Extract and list script steps from XML, validate generated FMXML snippets against 7 structural rules before they reach FileMaker, validate WebViewer HTML for remote dependencies and risky APIs, and write versioned artifacts to the project outbox for review.

### Specialist skills
Load curated FileMaker skill prompts (`pro_scriptwriter`, `script_analysis`, `fm_xml_serializer`, `script_debug`) directly into AI context to inject domain-specific guidance for writing, analyzing, or debugging scripts.

---

## Prerequisites

- **Exported FileMaker schema** — Use the KiBuild plugin's `Export Schema` button, or FileMaker's built-in DDR export tool, to produce an exploded schema folder:
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
- **macOS, Linux, or Windows** — Pre-compiled binaries are provided for all three.
- **No Go toolchain required** unless you are building from source.

### Where the schema XML can come from

The server indexes an **exploded** schema folder — one XML file *per object*, grouped into `scripts/`, `layouts/`, `tables/`, and `relationships/` (shown above). Several things can produce FileMaker XML; they differ in how much "exploding" is left to do:

| Source | Output | Ready to index? |
|---|---|---|
| KiBuild plugin **Export Schema** | Exploded tree, one file per object | ✅ Directly |
| Built-in **DDR** export | One `FMPReport` document for the whole solution | Not yet supported |
| **Save a Copy as XML** — single file | One `FMSaveAsXML` document for the whole solution | ✅ Via `explode_xml_export` |
| **Save a Copy as XML** — per-catalog option | One file per *catalog* (`<DB>_ScriptCatalog.xml`, `<DB>_LayoutCatalog.xml`, `<DB>_BaseTableCatalog.xml`, …) | ✅ Via `explode_xml_export` |

FileMaker 2025/2026's native **Save a Copy as XML** (available as both a menu command and a script step, configured via JSON options) is a convenient, license-friendly way to get schema XML out without the plugin. It emits the `FMSaveAsXML` dialect either as one document, or — with the `per catalog` option — as one `<DB>_<Catalog>Catalog.xml` file per catalog in a destination folder.

The catch: **per-catalog is not per-object.** Even with the split enabled, *all* scripts stay combined inside a single `ScriptCatalog`, all layouts inside `LayoutCatalog`, and so on — FileMaker does not explode each script into its own file, while the indexing tools here expect one file per object.

The built-in **`explode_xml_export`** tool closes that gap. Point it at either form (the single `FMSaveAsXML` file or the split-catalog folder — it auto-detects) and it writes the per-object layout the indexer needs, then run `generate_schema_map`:

```
Explode the Save-as-XML export at /path/to/Contacts.xml into my project, then build the schema map.
```

> **Coverage:** `explode_xml_export` explodes **every catalog** in the export, one file per object under `Schema/<database>/`:
> - `scripts/*.xml` + `scripts_sanitized/*.txt` — full script analysis and reference tools
> - `tables/*.xml` — base tables with their fields joined in (find_table, calculation/field references)
> - `layouts/*.xml`, `relationships/*.xml`, `table_occurrences/*.xml` — navigation, impact analysis, relationship predicates
> - `valuelists/`, `custom_functions/`, `custom_menus/`, `accounts/`, `privilege_sets/`, `extended_privileges/`, `themes/`, `base_directories/`, `external_data_sources/` — exploded for completeness and Git diffing
>
> The indexer reads `tables/`, `layouts/`, `relationships/`, `table_occurrences/`, and `scripts/`; the remaining folders are split out for browsing/version control. Script folders are flattened and name collisions disambiguated by id; relationship files are named after the joined table occurrences.

---

## Installation

### Step 1 — Get the binary

**Option A — Quick installation script (recommended)**

Run the install command for your platform:

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | sh
```
*Installs to `/usr/local/bin/kibuild-mcp`.*

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.ps1 | iex
```
*Installs to `%LOCALAPPDATA%\Programs\kibuild-mcp\kibuild-mcp.exe` and appends it to your User PATH.*

> **PowerShell execution policy:** If you get a "running scripts is disabled" error, run this first (one-time, per user):
> ```powershell
> Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
> ```

---

**Option B — Manual download**

Go to the [Releases page](https://github.com/priyabratasahoo21/kibuild-mcp/releases) and download the binary for your platform:

| Platform | File |
|---|---|
| macOS (Apple Silicon) | `kibuild-mcp-darwin-arm64` |
| macOS (Intel) | `kibuild-mcp-darwin-amd64` |
| Linux | `kibuild-mcp-linux-amd64` |
| Windows | `kibuild-mcp-windows-amd64.exe` |

Move it to a folder on your PATH:

- **macOS / Linux:**
  ```bash
  chmod +x kibuild-mcp-darwin-arm64
  mv kibuild-mcp-darwin-arm64 /usr/local/bin/kibuild-mcp
  ```
- **Windows (PowerShell):**
  ```powershell
  mkdir "$env:LOCALAPPDATA\Programs\kibuild-mcp"
  Move-Item kibuild-mcp-windows-amd64.exe "$env:LOCALAPPDATA\Programs\kibuild-mcp\kibuild-mcp.exe"
  [Environment]::SetEnvironmentVariable("PATH", $env:PATH + ";$env:LOCALAPPDATA\Programs\kibuild-mcp", "User")
  ```

> On macOS, the first time you run it you may need to allow it in **System Settings → Privacy & Security**.

---

**Option C — Go install**

Requires Go 1.21 or later.

```bash
go install github.com/priyabratasahoo21/kibuild-mcp@latest
```
*This places `kibuild-mcp` in your `$GOPATH/bin` (typically `~/go/bin`), which should be in your PATH.*

---

**Option D — Build from source**

Requires Go 1.21 or later.

```bash
git clone https://github.com/priyabratasahoo21/kibuild-mcp.git
cd kibuild-mcp
go build -ldflags="-s -w" -o kibuild-mcp .
mv kibuild-mcp /usr/local/bin/
```

---

## Setup

### Step 2 — Register in your AI tool

Pick your tool below. Paste the snippet into the config file shown, replacing `/path/to/your/project` with the absolute path to your FileMaker project folder.

> **Where is my project folder?** It's the folder that contains (or will contain) your `files/Schema/` export. The same path you would pass to `generate_schema_map`.

---

#### Claude Code

**macOS / Linux** — Config file: `~/.claude.json`

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

**Windows** — Config file: `C:\Users\<YourName>\.claude.json`

Open it in any text editor (create it if it doesn't exist) and use forward slashes or escaped backslashes for the path:

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

> **Windows tip:** Forward slashes (`/`) work in JSON on Windows. Double-backslashes (`C:\\Users\\...`) also work, but forward slashes are easier to read and less error-prone.

After editing, restart Claude Code (or run `/mcp` to reload servers).

---

#### OpenAI Codex CLI

Config file: `~/.codex/config.toml` (global) or `.codex/config.toml` in your repo (project-scoped)

```toml
[mcp_servers.kibuild]
command = "/usr/local/bin/kibuild-mcp"

[mcp_servers.kibuild.env]
KIBUILD_ACTIVE_PROJECT = "/path/to/your/project"
```

> **Windows path:** use `command = 'C:\Users\<YourName>\AppData\Local\Programs\kibuild-mcp\kibuild-mcp.exe'`

Or add it via the CLI:
```bash
codex mcp add kibuild -- /usr/local/bin/kibuild-mcp
```
Then open `~/.codex/config.toml` and add the `KIBUILD_ACTIVE_PROJECT` env var manually.

---

#### Google Antigravity (Agy)

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

Create the folder if it doesn't exist: `mkdir -p ~/.gemini/config`

This config is shared across Agy CLI, Antigravity IDE, and all other Antigravity tools.

---

#### Cursor

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

#### Windsurf

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

#### VS Code (with MCP extension)

Config file: User `settings.json` (`Ctrl+Shift+P` → "Open User Settings JSON")

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

### Step 3 — Build the workspace index

Once the server is registered, ask your AI tool:

**macOS / Linux:**
```
Call generate_schema_map for my project at /Users/yourname/Documents/MyFileMakerProject
```

**Windows:**
```
Call generate_schema_map for my project at C:/Users/YourName/Documents/MyFileMakerProject
```

> Use the same path you set in `KIBUILD_ACTIVE_PROJECT`. Forward slashes work on all platforms.

This scans your schema folder and writes `workspace_map.md` to your project root. After that, all navigation and reference tools are live. The index auto-refreshes whenever schema files change.

---

## Quick start

Once registered, ask your AI tool natural questions:

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

## Tool reference

### Schema navigation

| Tool | Description |
|---|---|
| `find_script` | Find a script by name. Returns sanitized step list, `txt_path`, `xml_path`, and sibling scripts in the same folder. Always call this before reading or modifying a script. |
| `find_layout` | Find a layout by name. Returns bound table occurrence, referenced scripts and layouts, and the layout XML path. |
| `find_table` | Find a base table by name. Returns all fields with types and the table XML path. |
| `inspect_relationships` | Return all relationship predicates for a database or table occurrence. |
| `search_index` | Keyword search over `workspace_map.md`. Returns only matching lines — token-efficient. Call `generate_schema_map` first if the index does not exist. |
| `generate_schema_map` | Build or refresh `workspace_map.md` — a compact Markdown index of all tables, layouts, scripts, and table occurrences across the workspace. |

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
| `xml_extract_steps` | List all script steps from a raw FMXML snippet or file content. |
| `xml_lookup_name` | Resolve a numeric script ID to its name from an XML document. |
| `xml_trace_dependencies` | Extract all referenced table occurrences, scripts, layouts, and fields from XML content. |
| `xml_match_revision` | Read the FileMaker version and revision metadata from an XML header. |
| `validate_fmxmlsnippet` | Run 7-rule structural validation on a generated FMXML snippet and return a pass/fail report with details. |
| `validate_webviewer_html` | Check generated WebViewer HTML for remote dependencies, risky JavaScript APIs, FileMaker bridge usage, and bundle size. |
| `write_outbox_artifact` | Save a generated script, layout, or document to the project outbox as a versioned artifact with a manifest entry. |
| `explode_xml_export` | Explode a FileMaker **Save a Copy as XML** export (single `FMSaveAsXML` file or split-catalog folder — auto-detected) into the per-object schema layout. Explodes every catalog — scripts (+ sanitized `.txt`), tables (fields joined in), layouts, relationships, table occurrences, value lists, custom functions/menus, accounts, privileges, themes — one file per object under `Schema/<database>/`; run `generate_schema_map` afterward. |

### Specialist skills

| Tool | Description |
|---|---|
| `load_skill` | Load a specialist skill by ID into AI context. Available skills: `pro_scriptwriter` (FileMaker scripting patterns), `script_analysis` (structured script audit), `fm_xml_serializer` (valid FMXML generation rules), `script_debug` (systematic debugging approach). |

---

## Logging

The server logs all MCP traffic to `~/.fm_ai_bridge/mcp_server.log`. Tail it to debug tool calls:

**macOS / Linux:**
```bash
tail -f ~/.fm_ai_bridge/mcp_server.log
```

**Windows (PowerShell):**
```powershell
Get-Content -Wait "$env:USERPROFILE\.fm_ai_bridge\mcp_server.log"
```

---

## Architecture

```
AI tool (Claude Code, Cursor, Windsurf, VS Code)
  │
  │  spawns subprocess on MCP connect
  ▼
kibuild-mcp  ← this binary
  │  MCP JSON-RPC over stdin/stdout (protocol 2024-11-05)
  │
  └── analysis tools  ← read Schema/ XML files on disk
        works without FileMaker running

Reads from disk:
  ~/your-project/files/Schema/<DBName>/   ← exported schema (XML files)
  ~/.fm_ai_bridge/active_project.txt      ← current project pointer
  ~/your-project/workspace_map.md         ← built on first run of generate_schema_map
```

---

## Contributing

Pull requests are welcome. Please open an issue first for significant changes.

---

## License

MIT

# KiBuild MCP

A self-contained MCP server that gives any AI coding tool deep, read-level access to a Claris FileMaker schema — no FileMaker license required at runtime for schema analysis.

Register one binary in your MCP config. Your AI tool (Claude Code, Cursor, Windsurf, VS Code) immediately gains FileMaker-aware tools: script navigation, impact analysis across the full dependency graph, XML generation and validation, specialist skills, and live FileMaker execution when the KiBuild plugin is connected.

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

### Live FileMaker execution (optional)
When the KiBuild C++ plugin is running and connected, `run_script` and `execute_sql` execute against the active FileMaker database over IPC. All other 34 tools work without FileMaker running.

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

---

## Installation

### Option 1 — Download a release binary (recommended)

1. Go to the [Releases page](https://github.com/priyabratasahoo21/kibuild-mcp/releases) and download the binary for your platform:

   | Platform | File |
   |---|---|
   | macOS (Apple Silicon) | `kibuild-mcp-darwin-arm64` |
   | macOS (Intel) | `kibuild-mcp-darwin-amd64` |
   | Linux | `kibuild-mcp-linux-amd64` |
   | Windows | `kibuild-mcp-windows-amd64.exe` |

2. Make it executable and move it to your PATH:

   ```bash
   # macOS / Linux
   chmod +x kibuild-mcp-darwin-arm64
   mv kibuild-mcp-darwin-arm64 /usr/local/bin/kibuild-mcp
   ```

   On macOS you may need to allow the binary in **System Settings → Privacy & Security** the first time you run it.

### Option 2 — Homebrew (macOS / Linux)

```bash
brew install kibuild/tap/kibuild-mcp
```

> Homebrew tap coming soon. Watch the repository for the release.

### Option 3 — Build from source

Requires Go 1.21 or later.

```bash
git clone https://github.com/priyabratasahoo21/kibuild-mcp.git
cd kibuild-mcp
go build -o kibuild-mcp .
mv kibuild-mcp /usr/local/bin/
```

### Option 4 — Install script

```bash
curl -fsSL https://raw.githubusercontent.com/priyabratasahoo21/kibuild-mcp/main/install.sh | bash
```

> Install script coming soon.

---

## Setup

### Step 1 — Point the binary at your project

**Option A — Environment variable** (recommended for MCP configs):

```json
"env": { "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project" }
```

**Option B — Active project file**:

```bash
mkdir -p ~/.fm_ai_bridge
echo "/path/to/your/project" > ~/.fm_ai_bridge/active_project.txt
```

### Step 2 — Register in your AI tool's MCP config

#### Claude Code

Add to `~/.claude.json`:

```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "args": [],
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

#### Cursor

Add to `~/.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "args": [],
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

#### Windsurf

Add to `~/.codeium/windsurf/mcp_config.json`:

```json
{
  "mcpServers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "args": [],
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

#### VS Code (with MCP extension)

Add to your VS Code `settings.json`:

```json
{
  "mcp.servers": {
    "kibuild": {
      "command": "/usr/local/bin/kibuild-mcp",
      "args": [],
      "env": {
        "KIBUILD_ACTIVE_PROJECT": "/path/to/your/project"
      }
    }
  }
}
```

### Step 3 — Build the workspace index

Ask your AI tool to call `generate_schema_map`. This scans the schema folder and writes `workspace_map.md` to your project root — after that, `search_index` and all navigation tools are live.

```
Call generate_schema_map for my project at /path/to/your/project
```

The index auto-refreshes whenever schema files change.

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

## Disabling specific tools

Add a `kibuild_config.json` file to `~/.fm_ai_bridge/`:

```json
{
  "disabled_mcp_tools": ["run_script", "execute_sql"]
}
```

Any tool name listed there is excluded from `tools/list` and blocked at call time.

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

### Specialist skills

| Tool | Description |
|---|---|
| `load_skill` | Load a specialist skill by ID into AI context. Available skills: `pro_scriptwriter` (FileMaker scripting patterns), `script_analysis` (structured script audit), `fm_xml_serializer` (valid FMXML generation rules), `script_debug` (systematic debugging approach). |

### Live FileMaker execution — requires KiBuild plugin

These tools require FileMaker to be open with the KiBuild plugin loaded. They gracefully error when the plugin is absent without affecting the other 34 tools.

| Tool | Description |
|---|---|
| `run_script` | Run a FileMaker script by name, with an optional parameter. |
| `execute_sql` | Execute an `ExecuteSQL` query against the active database and return results. |

---

## Logging

The server logs all MCP traffic to `~/.fm_ai_bridge/mcp_server.log`. Tail it to debug tool calls:

```bash
tail -f ~/.fm_ai_bridge/mcp_server.log
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
  ├── analysis tools  ← read Schema/ XML files on disk
  │     works with no FileMaker running
  │
  └──  2 live-exec tools  ← IPC → KiBuild C++ Plugin → FileMaker
        requires FileMaker open + KiBuild plugin loaded

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

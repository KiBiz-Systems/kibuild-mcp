---
description: Generate CLAUDE.md, AGENTS.md, GEMINI.md, and .cursor/rules/kibuild.mdc in the current project directory, instructing AI tools to use KiBuild MCP tools as the priority over grep/bash when working with FileMaker schemas.
---

You are initializing a FileMaker project to use KiBuild MCP tools as the primary AI interface. Generate all four guide files below in the **current working directory**.

If any file already exists, warn the user and ask before overwriting. Otherwise write all files automatically and report what was created.

---

## File 1 — CLAUDE.md (Claude Code)

Write this content to `CLAUDE.md`:

```markdown
# KiBuild MCP — Project AI Guide

This project uses the KiBuild MCP server for FileMaker-aware AI assistance.
**Always prefer KiBuild MCP tools over generic Bash, grep, or find.**

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script | `mcp__kibuild__find_script` |
| `grep` for a table | `mcp__kibuild__find_table` |
| `grep` for a layout | `mcp__kibuild__find_layout` |
| `grep` / `find` for keywords | `mcp__kibuild__search_index` |
| Reading XML files directly | `mcp__kibuild__xml_extract_steps` |
| `grep` for a name in XML | `mcp__kibuild__xml_lookup_name` |
| Tracing relationships manually | `mcp__kibuild__inspect_relationships` |
| `grep` for field usage | `mcp__kibuild__find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` / `_in_relationships` |
| `grep` for script call sites | `mcp__kibuild__find_script_references_in_scripts` / `_in_layouts` |
| `grep` for layout references | `mcp__kibuild__find_layout_references_in_calculations` / `_to_scripts` / `_to_tables` / `_to_valuelists` |
| `grep` for variable usage | `mcp__kibuild__find_variable_references_in_scripts` |
| Running a FileMaker SQL query | `mcp__kibuild__execute_sql` |
| Writing generated output | `mcp__kibuild__write_outbox_artifact` |

## Start every FileMaker session with

1. `mcp__kibuild__load_skill` — load the specialist skill for the task
2. `mcp__kibuild__search_index` — orient in the schema before diving in
3. Never `grep` or `cat` raw XML — use `xml_extract_steps` instead

## Other tools

- `mcp__kibuild__generate_schema_map` — rebuild index after an XML export
- `mcp__kibuild__list_workflows` / `mcp__kibuild__get_workflow` — available workflows
- `mcp__kibuild__validate_fmxmlsnippet` — validate generated FileMaker XML
- `mcp__kibuild__propose_preview` — propose a UI/layout change with preview
- `mcp__kibuild__run_script` — execute a FileMaker script via plugin
```

---

## File 2 — AGENTS.md (OpenAI Codex CLI)

Write this content to `AGENTS.md`:

```markdown
# KiBuild MCP — Project AI Guide (OpenAI Codex)

This project uses the KiBuild MCP server for FileMaker-aware AI assistance.
**Always prefer KiBuild MCP tools over shell commands like grep, find, or cat.**

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script | `kibuild__find_script` |
| `grep` for a table | `kibuild__find_table` |
| `grep` for a layout | `kibuild__find_layout` |
| `grep` / `find` for keywords | `kibuild__search_index` |
| Reading XML files directly | `kibuild__xml_extract_steps` |
| Tracing relationships manually | `kibuild__inspect_relationships` |
| `grep` for field usage | `kibuild__find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` / `_in_relationships` |
| `grep` for script call sites | `kibuild__find_script_references_in_scripts` / `_in_layouts` |
| `grep` for variable usage | `kibuild__find_variable_references_in_scripts` |
| Running a FileMaker SQL query | `kibuild__execute_sql` |
| Writing generated output | `kibuild__write_outbox_artifact` |

## Start every FileMaker session with

1. `kibuild__load_skill` — load the specialist skill
2. `kibuild__search_index` — orient in the schema
3. Never grep raw XML — use `kibuild__xml_extract_steps`
```

---

## File 3 — GEMINI.md (Google Gemini CLI)

Write this content to `GEMINI.md`:

```markdown
# KiBuild MCP — Project AI Guide (Google Gemini CLI)

This project uses the KiBuild MCP server for FileMaker-aware AI assistance.
**Always prefer KiBuild MCP tools over shell commands like grep, find, or cat.**

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script | `kibuild_find_script` |
| `grep` for a table | `kibuild_find_table` |
| `grep` for a layout | `kibuild_find_layout` |
| `grep` / `find` for keywords | `kibuild_search_index` |
| Reading XML files directly | `kibuild_xml_extract_steps` |
| Tracing relationships manually | `kibuild_inspect_relationships` |
| `grep` for field usage | `kibuild_find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` / `_in_relationships` |
| `grep` for script call sites | `kibuild_find_script_references_in_scripts` / `_in_layouts` |
| `grep` for variable usage | `kibuild_find_variable_references_in_scripts` |
| Running a FileMaker SQL query | `kibuild_execute_sql` |
| Writing generated output | `kibuild_write_outbox_artifact` |

## Start every FileMaker session with

1. `kibuild_load_skill` — load the specialist skill
2. `kibuild_search_index` — orient in the schema
3. Never grep raw XML — use `kibuild_xml_extract_steps`
```

---

## File 4 — .cursor/rules/kibuild.mdc (Cursor)

Create the directory `.cursor/rules/` if it does not exist, then write this content to `.cursor/rules/kibuild.mdc`:

```
---
description: KiBuild MCP tool preference rules for FileMaker projects
globs: "**"
alwaysApply: true
---

# KiBuild MCP Rules

This project uses the KiBuild MCP server. Always prefer KiBuild tools over grep/bash.

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script | `find_script` |
| `grep` for a table | `find_table` |
| `grep` for a layout | `find_layout` |
| `grep` / `find` for keywords | `search_index` |
| Reading XML files directly | `xml_extract_steps` |
| Tracing relationships manually | `inspect_relationships` |
| `grep` for field usage | `find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` |
| `grep` for script call sites | `find_script_references_in_scripts` / `_in_layouts` |
| `grep` for variable usage | `find_variable_references_in_scripts` |
| Running FileMaker SQL | `execute_sql` |
| Writing generated output | `write_outbox_artifact` |

## Session startup

1. Call `load_skill` first — loads the right specialist context
2. Call `search_index` to orient before diving in
3. Never cat or grep raw XML — use `xml_extract_steps`
```

---

## After writing all files

Report:

```
✓ CLAUDE.md              — Claude Code guide
✓ AGENTS.md              — OpenAI Codex guide
✓ GEMINI.md              — Google Gemini CLI guide
✓ .cursor/rules/kibuild.mdc — Cursor guide

All files instruct the AI to prefer KiBuild MCP tools over grep/bash.
Commit these files to your repo so all team members get the same behavior.

To regenerate: /init-kibuild-project
```

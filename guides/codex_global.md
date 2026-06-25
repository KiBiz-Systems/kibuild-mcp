# KiBuild MCP — Global Instructions (OpenAI Codex CLI)

This machine has KiBuild MCP installed for FileMaker-aware AI assistance.
When working in a FileMaker project, **always prefer KiBuild MCP tools over shell commands like grep, find, or cat on XML files.**

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script name | `kibuild__find_script` |
| `grep` for a table name | `kibuild__find_table` |
| `grep` for a layout name | `kibuild__find_layout` |
| `grep` / `find` for keywords | `kibuild__search_index` |
| Reading or catting XML files | `kibuild__xml_extract_steps` |
| `grep` for a name in XML | `kibuild__xml_lookup_name` |
| Tracing relationships manually | `kibuild__inspect_relationships` |
| `grep` for where a field is used | `kibuild__find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` / `_in_relationships` |
| `grep` for script call sites | `kibuild__find_script_references_in_scripts` / `_in_layouts` |
| `grep` for layout references | `kibuild__find_layout_references_in_calculations` / `_to_scripts` |
| `grep` for variable usage | `kibuild__find_variable_references_in_scripts` |
| Running a FileMaker SQL query | `kibuild__execute_sql` |
| Writing generated FileMaker output | `kibuild__write_outbox_artifact` |

## Start every FileMaker session with

1. `kibuild__load_skill` — loads the right specialist context for the task
2. `kibuild__search_index` — orient in the schema before diving into specifics
3. Never grep or cat raw XML — use `kibuild__xml_extract_steps` instead

## Other useful tools

- `kibuild__generate_schema_map` — rebuild the schema index after an XML export
- `kibuild__list_workflows` / `kibuild__get_workflow` — list and inspect workflows
- `kibuild__validate_fmxmlsnippet` — validate generated FileMaker XML before writing
- `kibuild__run_script` — execute a FileMaker script via the plugin

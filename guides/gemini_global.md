# KiBuild MCP — Global Instructions (Google Gemini CLI)

This machine has KiBuild MCP installed for FileMaker-aware AI assistance.
When working in a FileMaker project, **always prefer KiBuild MCP tools over shell commands like grep, find, or cat on XML files.**

## Tool substitution rules

| Instead of... | Use this KiBuild MCP tool |
|---|---|
| `grep` for a script name | `kibuild_find_script` |
| `grep` for a table name | `kibuild_find_table` |
| `grep` for a layout name | `kibuild_find_layout` |
| `grep` / `find` for keywords | `kibuild_search_index` |
| Reading or catting XML files | `kibuild_xml_extract_steps` |
| `grep` for a name in XML | `kibuild_xml_lookup_name` |
| Tracing relationships manually | `kibuild_inspect_relationships` |
| `grep` for where a field is used | `kibuild_find_field_references_in_calculations` / `_in_scripts` / `_in_layouts` / `_in_relationships` |
| `grep` for script call sites | `kibuild_find_script_references_in_scripts` / `_in_layouts` |
| `grep` for layout references | `kibuild_find_layout_references_in_calculations` / `_to_scripts` |
| `grep` for variable usage | `kibuild_find_variable_references_in_scripts` |
| Running a FileMaker SQL query | `kibuild_execute_sql` |
| Writing generated FileMaker output | `kibuild_write_outbox_artifact` |

## Start every FileMaker session with

1. `kibuild_load_skill` — loads the right specialist context for the task
2. `kibuild_search_index` — orient in the schema before diving into specifics
3. Never grep or cat raw XML — use `kibuild_xml_extract_steps` instead

## Other useful tools

- `kibuild_generate_schema_map` — rebuild the schema index after an XML export
- `kibuild_list_workflows` / `kibuild_get_workflow` — list and inspect workflows
- `kibuild_validate_fmxmlsnippet` — validate generated FileMaker XML before writing
- `kibuild_run_script` — execute a FileMaker script via the plugin

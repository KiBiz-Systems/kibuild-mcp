package tools

import (
	"encoding/json"
	"github.com/priyabratasahoo21/kibuild-mcp/providers"
)

func GetToolsSchema() []providers.Tool {
	allTools := []providers.Tool{
		{
			Name:        "list_dir",
			Description: "List contents of a directory. Returns names, sizes and whether they are directories.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"path": {
						"type": "string",
						"description": "Absolute path to the directory"
					}
				},
				"required": ["path"]
			}`),
		},
		{
			Name:        "read_file",
			Description: "Read the full contents of a file at the specified path.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"path": {
						"type": "string",
						"description": "Absolute path to the file"
					}
				},
				"required": ["path"]
			}`),
		},
		{
			Name:        "write_file",
			Description: "Write contents to a file at the specified path. Overwrites the file if it exists.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"path": {
						"type": "string",
						"description": "Absolute path to the file"
					},
					"content": {
						"type": "string",
						"description": "The exact content to write"
					}
				},
				"required": ["path", "content"]
			}`),
		},
		{
			Name:        "run_command",
			Description: "Propose a terminal command to run. Safe shell execution.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"command": {
						"type": "string",
						"description": "The command line string to run"
					}
				},
				"required": ["command"]
			}`),
		},
		{
			Name:        "export_schema",
			Description: "Export the current FileMaker schema using the exploder binary.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"database": {
						"type": "string",
						"description": "The name of the database file (e.g. Sales.fmp12)"
					}
				},
				"required": ["database"]
			}`),
		},
		{
			Name:        "read_layout",
			Description: "Read the XML/JSON representation of a FileMaker layout.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"layout_name": {
						"type": "string",
						"description": "Name of the layout to read"
					}
				},
				"required": ["layout_name"]
			}`),
		},
		{
			Name:        "get_active_context",
			Description: "Get context about the active FileMaker database, layout name, and currently selected record.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {}
			}`),
		},
		{
			Name:        "xml_extract_steps",
			Description: "Extract and list script steps from a FileMaker XML script snippet or file.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"xml_content": {
						"type": "string",
						"description": "The FileMaker XML content of the script"
					}
				},
				"required": ["xml_content"]
			}`),
		},
		{
			Name:        "xml_lookup_name",
			Description: "Lookup a script name by its ID in the FileMaker XML schemas.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"xml_content": {
						"type": "string",
						"description": "The XML content containing script references"
					},
					"script_id": {
						"type": "string",
						"description": "The ID of the script to lookup"
					}
				},
				"required": ["xml_content", "script_id"]
			}`),
		},
		{
			Name:        "xml_trace_dependencies",
			Description: "Trace layout or script dependencies (referenced table occurrences, scripts, layouts, fields) in the XML content.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"xml_content": {
						"type": "string",
						"description": "The XML content to trace"
					}
				},
				"required": ["xml_content"]
			}`),
		},
		{
			Name:        "xml_match_revision",
			Description: "Compare or extract FileMaker version and revision parameters from the XML header metadata.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"xml_content": {
						"type": "string",
						"description": "The XML content containing root metadata"
					}
				},
				"required": ["xml_content"]
			}`),
		},
		{
			Name:        "diff_patch",
			Description: "Overwrites file content at the target path with patched content.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"path": {
						"type": "string",
						"description": "The absolute path of the file to modify"
					},
					"content": {
						"type": "string",
						"description": "The new content to write to the file"
					}
				},
				"required": ["path", "content"]
			}`),
		},
		{
			Name:        "generate_schema_map",
			Description: "Generate a compact markdown map (workspace_map.md) of the workspace tables, fields, scripts, and layouts, facilitating RAG analysis.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"project_path": {
						"type": "string",
						"description": "The absolute path of the project workspace to scan"
					}
				},
				"required": ["project_path"]
			}`),
		},
		{
			Name:        "search_file",
			Description: "Search for a text pattern in files under the given directory (like grep -rn). Returns matching file names, line numbers, and content. Capped at 50 matches.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"pattern": {
						"type": "string",
						"description": "The text pattern to search for"
					},
					"path": {
						"type": "string",
						"description": "The directory path to search in"
					},
					"case_insensitive": {
						"type": "boolean",
						"description": "If true, performs case-insensitive search. Defaults to false."
					}
				},
				"required": ["pattern", "path"]
			}`),
		},
		{
			Name:        "read_xml_guide",
			Description: "Read the FileMaker XML snippet reference guide containing schemas and XML templates for script step generation.",
			Parameters:  json.RawMessage(`{"type":"object","properties":{}}`),
		},
		{
			Name:        "validate_fmxmlsnippet",
			Description: "Validate a FileMaker XML snippet against the structural rules and return a validation report.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"xml_content": {
						"type": "string",
						"description": "The XML snippet to validate"
					},
					"database": {
						"type": "string",
						"description": "Optional name of the target database file for context checks (e.g. Sales.fmp12)"
					}
				},
				"required": ["xml_content"]
			}`),
		},
		{
			Name:        "write_outbox_artifact",
			Description: "Write a generated script, layout, or document to the project outbox, creating a versioned folder and tracking it in the manifest.json. This is preferred over standard write_file for generated outputs.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"artifact_id": {
						"type": "string",
						"description": "Unique slug for the artifact (e.g. create_contact)"
					},
					"artifact_type": {
						"type": "string",
						"enum": ["script", "layout", "schema", "doc"],
						"description": "The type of generated artifact"
					},
					"artifact_name": {
						"type": "string",
						"description": "Human-friendly display name (e.g. Create Contact)"
					},
					"database": {
						"type": "string",
						"description": "Name of the target database file"
					},
					"files": {
						"type": "object",
						"description": "Map of filename to exact file content. For scripts, provide the compiled .json payload and a sanitized .txt representation.",
						"additionalProperties": {
							"type": "string"
						}
					}
				},
				"required": ["artifact_id", "artifact_type", "artifact_name", "database", "files"]
			}`),
		},
		{
			Name:        "find_script",
			Description: "Find a FileMaker script by name in the project Schema directory. Returns: sanitized step list (.txt), txt_path, xml_path, and sibling scripts in the same folder. ALWAYS call this before reading or modifying any script. Use txt_path content for analysis; only read xml_path when generating output XML.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"script_name": {
						"type": "string",
						"description": "Name or partial name of the script to find (case-insensitive fuzzy match)"
					},
					"database": {
						"type": "string",
						"description": "Optional database name (e.g. Contacts) to narrow the search. Omit to search all databases."
					}
				},
				"required": ["script_name"]
			}`),
		},
		{
			Name:        "find_table",
			Description: "Find a FileMaker base table by name in exploded schema and return fields, table XML path, and database context. Use this before proposing table/field changes.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"table_name": {
						"type": "string",
						"description": "Name or partial name of the base table to find"
					},
					"database": {
						"type": "string",
						"description": "Optional database name to narrow the search"
					}
				},
				"required": ["table_name"]
			}`),
		},
		{
			Name:        "find_layout",
			Description: "Find a FileMaker layout by name in exploded schema and return bound table occurrence, referenced scripts/layouts, and layout XML path.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"layout_name": {
						"type": "string",
						"description": "Name or partial name of the layout to find"
					},
					"database": {
						"type": "string",
						"description": "Optional database name to narrow the search"
					}
				},
				"required": ["layout_name"]
			}`),
		},
		{
			Name:        "inspect_relationships",
			Description: "Read exploded FileMaker relationship XML and return relationship predicates for a database or table occurrence.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"database": {
						"type": "string",
						"description": "Optional database name to inspect"
					},
					"table_occurrence": {
						"type": "string",
						"description": "Optional table occurrence name to filter relationships"
					}
				}
			}`),
		},
		{
			Name:        "validate_webviewer_html",
			Description: "Validate generated WebViewer HTML for KiBuild preview/export. Checks remote dependencies, risky JavaScript APIs, FileMaker bridge usage, and bundle size.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"html": {
						"type": "string",
						"description": "Generated self-contained HTML intended for a FileMaker WebViewer"
					},
					"allow_remote_assets": {
						"type": "boolean",
						"description": "If true, remote scripts/styles/images are warnings instead of errors. Defaults to false."
					}
				},
				"required": ["html"]
			}`),
		},
		{
			Name:        "search_index",
			Description: "Search the workspace_map.md index for scripts, layouts, tables, or fields matching a keyword. Returns only matching lines (token-efficient). Use this INSTEAD of reading workspace_map.md directly. Call generate_schema_map first if the index does not exist.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"query": {
						"type": "string",
						"description": "Keyword to search for (case-insensitive). E.g. 'Daily Sales', 'REPORTS', 'Contact'"
					},
					"type": {
						"type": "string",
						"enum": ["script", "layout", "table", "field", "all"],
						"description": "Filter results by type. Use 'all' or omit to search everything."
					}
				},
				"required": ["query"]
			}`),
		},
		{
			Name:        "run_script",
			Description: "Run a FileMaker script by name, optionally passing a parameter.",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "script_name": {
				            "type": "string",
				            "description": "Name of the script to run"
				        },
				        "parameter": {
				            "type": "string",
				            "description": "Optional parameter to pass to the script"
				        }
				    },
				    "required": [
				        "script_name"
				    ]
				}`),
		},
		{
			Name:        "execute_sql",
			Description: "Execute a FileMaker ExecuteSQL query against the active database.",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "query": {
				            "type": "string",
				            "description": "The SQL query to execute"
				        }
				    },
				    "required": [
				        "query"
				    ]
				}`),
		},
		{
			Name:        "find_layout_references_to_scripts",
			Description: "Find references related to layout names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of layout names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_layout_references_to_valuelists",
			Description: "Find references related to layout names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of layout names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_layout_references_to_tables",
			Description: "Find references related to layout names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of layout names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_script_references_in_scripts",
			Description: "Find references related to script names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of script names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_script_references_in_layouts",
			Description: "Find references related to script names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of script names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_script_references_to_layouts",
			Description: "Find references related to script names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of script names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_script_references_to_valuelists",
			Description: "Find references related to script names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of script names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_field_references_in_scripts",
			Description: "Find references related to field names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of field names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_field_references_in_layouts",
			Description: "Find references related to field names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of field names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_field_references_in_calculations",
			Description: "Find references related to field names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of field names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_field_references_in_relationships",
			Description: "Find references related to field names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of field names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_variable_references_in_scripts",
			Description: "Find references related to variable names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of variable names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_valuelist_references_in_calculations",
			Description: "Find references related to value list names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of value list names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_layout_references_in_calculations",
			Description: "Find references related to layout names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of layout names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_to_references",
			Description: "Find references related to table occurrence names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of table occurrence names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
		{
			Name:        "find_relationship_predicates",
			Description: "Find references related to table occurrence names",
			Parameters: json.RawMessage(`{
				    "type": "object",
				    "properties": {
				        "names": {
				            "type": "array",
				            "items": {
				                "type": "string"
				            },
				            "description": "List of table occurrence names to inspect"
				        },
				        "database": {
				            "type": "string",
				            "description": "Database name context"
				        }
				    },
				    "required": [
				        "names"
				    ]
				}`),
		},
	}

	allTools = append(allTools,
		providers.Tool{
			Name:        "explode_xml_export",
			Description: "Explode a FileMaker 'Save a Copy as XML' export into the per-object schema layout the navigation and reference tools index. Accepts either a single FMSaveAsXML file (split_catalogs=False) or a folder of split *Catalog.xml files (split_catalogs=True). Explodes every catalog: scripts (with sanitized .txt), tables (fields joined in), layouts, relationships, table occurrences, value lists, custom functions, custom menus, accounts, privilege sets, themes, and more — one file per object under <dest>/Schema/<database>/. Run generate_schema_map afterward to index the result.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"source": {
						"type": "string",
						"description": "Path to the FileMaker XML export: either a single FMSaveAsXML .xml file or a folder containing split per-catalog files (e.g. Contacts_ScriptCatalog.xml). Absolute, or relative to the active project."
					},
					"database": {
						"type": "string",
						"description": "Target database name for the exploded Schema/<database>/ folder. Inferred from the source file/folder name when omitted."
					},
					"dest": {
						"type": "string",
						"description": "Optional destination root. Defaults to <project>/files so output lands at <project>/files/Schema/<database>/."
					}
				},
				"required": ["source"]
			}`),
		},
		providers.Tool{
			Name:        "load_skill",
			Description: "Load the full instruction content of a specialist skill by its ID (e.g. 'pro_scriptwriter', 'script_analysis', 'fm_xml_serializer', 'script_debug'). Injects FileMaker-specific guidance into AI context for the current task.",
			Parameters: json.RawMessage(`{
				"type": "object",
				"properties": {
					"skill_id": {
						"type": "string",
						"description": "The skill ID to load (e.g. 'pro_scriptwriter', 'script_analysis', 'fm_xml_serializer', 'script_debug')"
					}
				},
				"required": ["skill_id"]
			}`),
		},
	)

	if !IsPluginConnected() {
		var filtered []providers.Tool
		for _, t := range allTools {
			if t.Name != "export_schema" && t.Name != "read_layout" && t.Name != "get_active_context" {
				filtered = append(filtered, t)
			}
		}
		return filtered
	}

	return allTools
}

var safeMCPTools = map[string]bool{
	"run_script":                                true,
	"execute_sql":                               true,
	"find_layout_references_to_scripts":         true,
	"find_layout_references_to_valuelists":      true,
	"find_layout_references_to_tables":          true,
	"find_script_references_in_scripts":         true,
	"find_script_references_in_layouts":         true,
	"find_script_references_to_layouts":         true,
	"find_script_references_to_valuelists":      true,
	"find_field_references_in_scripts":          true,
	"find_field_references_in_layouts":          true,
	"find_field_references_in_calculations":     true,
	"find_field_references_in_relationships":    true,
	"find_variable_references_in_scripts":       true,
	"find_valuelist_references_in_calculations": true,
	"find_layout_references_in_calculations":    true,
	"find_to_references":                        true,
	"find_relationship_predicates":              true,
	"export_schema":                             true,
	"read_layout":                               true,
	"search_index":                              true,
	"generate_schema_map":                       true,
	"find_script":                               true,
	"find_table":                                true,
	"find_layout":                               true,
	"inspect_relationships":                     true,
	"xml_trace_dependencies":                    true,
	"xml_extract_steps":                         true,
	"xml_lookup_name":                           true,
	"xml_match_revision":                        true,
	"validate_fmxmlsnippet":                     true,
	"validate_webviewer_html":                   true,
	"write_outbox_artifact":                     true,
	"explode_xml_export":                        true,
	"get_active_context":                        true,
	"read_xml_guide":                            true,
	"load_skill":                                true,
}

func IsSafeMCPTool(name string) bool {
	return safeMCPTools[name]
}

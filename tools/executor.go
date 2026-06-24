package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"kibuild/mcp/outbox"
	"kibuild/mcp/providers"
	"kibuild/mcp/validator"
)

func ExecuteTool(ctx context.Context, name string, argsJSON string) (string, error) {
	projectPath, _ := ctx.Value(providers.ProjectPathKey).(string)
	if projectPath == "" {
		if cwd, err := os.Getwd(); err == nil {
			projectPath = cwd
		}
	}

	switch name {
	case "list_dir":
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.Path)
		if err != nil {
			return "", err
		}
		return ListDir(sandboxed)

	case "read_file":
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.Path)
		if err != nil {
			return "", err
		}
		return ReadFile(sandboxed)

	case "write_file":
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.Path)
		if err != nil {
			return "", err
		}
		result, writeErr := WriteFile(sandboxed, args.Content)

		// Auto-generate .txt sanitized version when writing an fmxmlsnippet .xml to Outbox
		if writeErr == nil &&
			strings.HasSuffix(strings.ToLower(sandboxed), ".xml") &&
			strings.Contains(args.Content, "fmxmlsnippet") {
			if sanitized, sanitizeErr := SanitizeFMXmlSnippet(args.Content); sanitizeErr == nil {
				txtPath := sandboxed[:len(sandboxed)-4] + ".txt"
				header := fmt.Sprintf("# %s\n# Auto-generated sanitized script\n\n",
					strings.TrimSuffix(filepath.Base(sandboxed), ".xml"))
				_, _ = WriteFile(txtPath, header+sanitized)
				result += fmt.Sprintf(" | Auto-generated sanitized text: %s", filepath.Base(txtPath))
			}
		}
		return result, writeErr

	case "run_command":
		var args struct {
			Command string `json:"command"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return RunCommand(ctx, args.Command)

	case "export_schema":
		var args struct {
			Database string `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return ExportSchema(args.Database)

	case "read_layout":
		var args struct {
			LayoutName string `json:"layout_name"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return ReadLayout(args.LayoutName)

	case "get_active_context":
		return GetActiveContext()

	case "xml_extract_steps":
		var args struct {
			XMLContent string `json:"xml_content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return ExtractScriptSteps(args.XMLContent)

	case "xml_lookup_name":
		var args struct {
			XMLContent string `json:"xml_content"`
			ScriptID   string `json:"script_id"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return LookupScriptName(args.XMLContent, args.ScriptID)

	case "xml_trace_dependencies":
		var args struct {
			XMLContent string `json:"xml_content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return TraceDependencies(args.XMLContent)

	case "xml_match_revision":
		var args struct {
			XMLContent string `json:"xml_content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return MatchRevision(args.XMLContent)

	case "diff_patch":
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.Path)
		if err != nil {
			return "", err
		}
		return DiffPatch(sandboxed, args.Content)

	case "generate_schema_map":
		var args struct {
			ProjectPath string `json:"project_path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.ProjectPath)
		if err != nil {
			return "", err
		}
		return GenerateSchemaMap(sandboxed)

	case "search_file":
		var args struct {
			Pattern         string `json:"pattern"`
			Path            string `json:"path"`
			CaseInsensitive bool   `json:"case_insensitive"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		sandboxed, err := ResolveAndSandboxPath(projectPath, args.Path)
		if err != nil {
			return "", err
		}
		return SearchFile(sandboxed, args.Pattern, args.CaseInsensitive)

	case "read_xml_guide":
		return ReadXMLGuide(projectPath)

	case "validate_fmxmlsnippet":
		var args struct {
			XmlContent string `json:"xml_content"`
			Database   string `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		res, err := validator.ValidateXML(args.XmlContent, projectPath, args.Database)
		if err != nil {
			return "", err
		}
		resBytes, err := json.Marshal(res)
		if err != nil {
			return "", err
		}
		return string(resBytes), nil

	case "write_outbox_artifact":
		var args struct {
			ArtifactID   string            `json:"artifact_id"`
			ArtifactType string            `json:"artifact_type"`
			ArtifactName string            `json:"artifact_name"`
			Database     string            `json:"database"`
			Files        map[string]string `json:"files"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		if args.ArtifactID == "" || args.ArtifactType == "" || args.ArtifactName == "" || args.Database == "" || len(args.Files) == 0 {
			return fmt.Errorf("missing required fields in outbox artifact payload").Error(), nil // return clean error
		}

		// Emit immediate ack so the CLI agent doesn't interpret compile silence as a crash
		if logFn, ok := ctx.Value(providers.LogFuncKey).(func(string, ...interface{})); ok {
			logFn("[TOOL] write_outbox_artifact received for %q (%s) — compiling...", args.ArtifactName, args.ArtifactType)
		}

		// Intercept JSON payloads for scripts: compile to XML and generate .txt from the XML.
		if args.ArtifactType == "script" {
			for filename, content := range args.Files {
				if strings.HasSuffix(strings.ToLower(filename), ".json") {
					// Race compilation against a 30s timeout to prevent silent hangs on large payloads
					type compileResult struct {
						xml string
						err error
					}
					ch := make(chan compileResult, 1)
					go func(data []byte) {
						xml, err := CompileScript(projectPath, data)
						ch <- compileResult{xml, err}
					}([]byte(content))

					var compiledXml string
					var compileErr error
					select {
					case res := <-ch:
						compiledXml, compileErr = res.xml, res.err
					case <-time.After(30 * time.Second):
						return "script compilation timed out after 30s — the JSON payload may be too large or malformed. Split into smaller chunks and retry.", nil
					}

					if compileErr != nil {
						return fmt.Sprintf("script compilation failed for '%s': %v\nFix the step names or parameters and try again.", filename, compileErr), nil
					}
					if compiledXml != "" {
						baseName := strings.TrimSuffix(filename, filepath.Ext(filename))
						args.Files[baseName+".xml"] = compiledXml
						// Generate .txt from the compiled XML so both tabs show the same logic.
						// This overwrites any .txt the agent provided, which may differ from the JSON.
						if sanitized, sanitErr := SanitizeFMXmlSnippet(compiledXml); sanitErr == nil && sanitized != "" {
							args.Files[baseName+".txt"] = sanitized
						}
					}
				}
			}
		}

		// Ensure the artifact is registered first
		art, err := outbox.CreateArtifact(projectPath, args.ArtifactID, args.ArtifactType, args.ArtifactName, args.Database)
		if err != nil {
			return "", fmt.Errorf("failed to register outbox artifact: %w", err)
		}

		// Write version
		v, err := outbox.WriteVersion(projectPath, art.ID, args.Files)
		if err != nil {
			return "", fmt.Errorf("failed to write version for outbox: %w", err)
		}

		// Collect written file paths for response
		outboxPath := outbox.ResolveOutboxPath(projectPath)
		typeFolder := art.Type + "s"
		versionFolder := filepath.Join(outboxPath, typeFolder, art.ID, v.VersionID)

		var xmlAbsPath, txtAbsPath string
		for _, relFile := range v.Files {
			base := filepath.Base(relFile)
			if strings.HasSuffix(strings.ToLower(base), ".xml") {
				xmlAbsPath = filepath.Join(projectPath, relFile)
			} else if strings.HasSuffix(strings.ToLower(base), ".txt") {
				txtAbsPath = filepath.Join(projectPath, relFile)
			}
		}

		// Write diff.json — summary for the Changes tab in OutboxPage
		diffData := map[string]interface{}{
			"artifact_id": art.ID,
			"version":     v.VersionID,
			"script_name": art.Name,
			"database":    art.Database,
			"timestamp":   time.Now().Format(time.RFC3339),
			"files": map[string]string{
				"xml": xmlAbsPath,
				"txt": txtAbsPath,
			},
			"note": "Review changes in the Outbox panel. Accept or Reject this version.",
		}
		
		// If agent provided a diff.json, merge it
		if agentDiffStr, ok := args.Files["diff.json"]; ok {
			var agentDiff map[string]interface{}
			if err := json.Unmarshal([]byte(agentDiffStr), &agentDiff); err == nil {
				for k, val := range agentDiff {
					diffData[k] = val
				}
			}
		}

		if diffBytes, errMarshal := json.MarshalIndent(diffData, "", "  "); errMarshal == nil {
			diffPath := filepath.Join(versionFolder, "diff.json")
			_ = os.WriteFile(diffPath, diffBytes, 0644)
		}

		// Build response with clickable file:// links
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("✅ Saved %s (%s) to Outbox\n", art.Name, v.VersionID))
		if xmlAbsPath != "" {
			sb.WriteString(fmt.Sprintf("📄 XML: file://%s\n", xmlAbsPath))
		}
		if txtAbsPath != "" {
			sb.WriteString(fmt.Sprintf("📝 TXT: file://%s\n", txtAbsPath))
		}
		sb.WriteString("Open the Outbox panel to review, accept, or reject this version.")
		return sb.String(), nil

	case "find_script":
		var args struct {
			ScriptName string `json:"script_name"`
			Database   string `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindScript(projectPath, args.ScriptName, args.Database)

	case "find_table":
		var args struct {
			TableName string `json:"table_name"`
			Database  string `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindTable(projectPath, args.TableName, args.Database)

	case "find_layout":
		var args struct {
			LayoutName string `json:"layout_name"`
			Database   string `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindLayout(projectPath, args.LayoutName, args.Database)

	case "inspect_relationships":
		var args struct {
			Database        string `json:"database"`
			TableOccurrence string `json:"table_occurrence"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return InspectRelationships(projectPath, args.Database, args.TableOccurrence)

	case "validate_webviewer_html":
		var args struct {
			HTML              string `json:"html"`
			AllowRemoteAssets bool   `json:"allow_remote_assets"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return ValidateWebViewerHTML(args.HTML, args.AllowRemoteAssets)

	case "run_script":
		var args struct {
			ScriptName string `json:"script_name"`
			Parameter  string `json:"parameter"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return RunScript(args.ScriptName, args.Parameter)

	case "execute_sql":
		var args struct {
			Query string `json:"query"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return ExecuteSQL(args.Query)

	case "find_layout_references_to_scripts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindLayoutReferencesToScripts(projectPath, args.Names, args.Database)

	case "find_layout_references_to_valuelists":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindLayoutReferencesToValueLists(projectPath, args.Names, args.Database)

	case "find_layout_references_to_tables":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindLayoutReferencesToTables(projectPath, args.Names, args.Database)

	case "find_script_references_in_scripts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindScriptReferencesInScripts(projectPath, args.Names, args.Database)

	case "find_script_references_in_layouts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindScriptReferencesInLayouts(projectPath, args.Names, args.Database)

	case "find_script_references_to_layouts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindScriptReferencesToLayouts(projectPath, args.Names, args.Database)

	case "find_script_references_to_valuelists":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindScriptReferencesToValueLists(projectPath, args.Names, args.Database)

	case "find_field_references_in_scripts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindFieldReferencesInScripts(projectPath, args.Names, args.Database)

	case "find_field_references_in_layouts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindFieldReferencesInLayouts(projectPath, args.Names, args.Database)

	case "find_field_references_in_calculations":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindFieldReferencesInCalculations(projectPath, args.Names, args.Database)

	case "find_field_references_in_relationships":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindFieldReferencesInRelationships(projectPath, args.Names, args.Database)

	case "find_variable_references_in_scripts":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindVariableReferencesInScripts(projectPath, args.Names, args.Database)

	case "find_valuelist_references_in_calculations":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindValueListReferencesInCalculations(projectPath, args.Names, args.Database)

	case "find_layout_references_in_calculations":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindLayoutReferencesInCalculations(projectPath, args.Names, args.Database)

	case "find_to_references":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindTOReferences(projectPath, args.Names, args.Database)

	case "find_relationship_predicates":
		var args struct {
			Names    []string `json:"names"`
			Database string   `json:"database"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return FindRelationshipPredicates(projectPath, args.Names, args.Database)

	case "search_index":
		var args struct {
			Query      string `json:"query"`
			FilterType string `json:"type"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return SearchIndex(projectPath, args.Query, args.FilterType)

	case "load_skill":
		var args struct {
			SkillID string `json:"skill_id"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", fmt.Errorf("invalid arguments: %v", err)
		}
		return LoadSkillTool(args.SkillID)

	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

func RunCommand(ctx context.Context, cmdStr string) (string, error) {
	cmdCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("command failed: %v, output: %s", err, string(out))
	}
	return string(out), nil
}

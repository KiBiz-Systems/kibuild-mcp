// Package exploder converts FileMaker "Save a Copy as XML" output (the
// FMSaveAsXML dialect) into the exploded, one-file-per-object schema layout
// that the navigation and reference tools index.
//
// It accepts either form FileMaker produces:
//   - a single FMSaveAsXML file (split_catalogs="False") containing every
//     *Catalog under one <Structure><AddAction>, or
//   - a folder of split catalog files (split_catalogs="True"), one
//     <DB>_<Catalog>Catalog.xml per catalog.
//
// Both are the same dialect; the parsers stream and locate catalog elements by
// name regardless of nesting, so the single file is simply handed to every
// catalog parser while the split folder routes each file to its parser.
package exploder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Sanitizer turns an fmxmlsnippet into human-readable script text. It is
// injected by the caller (the tools package owns SanitizeFMXmlSnippet) so this
// package stays dependency-free and avoids an import cycle. May be nil.
type Sanitizer func(snippet string) (string, error)

// Result reports what was written.
type Result struct {
	Database string         `json:"database"`
	Dest     string         `json:"dest"`
	Source   string         `json:"source"`
	Mode     string         `json:"mode"`   // "single-file" or "split-catalogs"
	Counts   map[string]int `json:"counts"` // output folder -> objects written
	Total    int            `json:"total"`
	Warnings []string       `json:"warnings,omitempty"`
}

// Explode reads the FileMaker XML export at source and writes the exploded
// schema under dest/Schema/<database>/. If database is empty it is inferred
// from the file/folder name. dest defaults to the source's parent when empty.
func Explode(source, database, dest string, sanitize Sanitizer) (*Result, error) {
	info, err := os.Stat(source)
	if err != nil {
		return nil, fmt.Errorf("source not found: %w", err)
	}
	isDir := info.IsDir()

	if database == "" {
		database = inferDatabase(source, isDir)
	}
	if dest == "" {
		dest = filepath.Dir(source)
	}

	res := &Result{Database: database, Source: source, Counts: map[string]int{}}
	if isDir {
		res.Mode = "split-catalogs"
	} else {
		res.Mode = "single-file"
	}

	schemaRoot := filepath.Join(dest, "Schema", database)
	res.Dest = schemaRoot

	record := func(folder string, n int, err error) {
		if err != nil {
			res.Warnings = append(res.Warnings, fmt.Sprintf("%s: %v", folder, err))
		}
		if n > 0 {
			res.Counts[folder] = n
			res.Total += n
		}
	}

	// Scripts: the step list lives in <StepsForScripts>, sanitized to .txt.
	if xmlStr, ok := loadCatalogXML(source, isDir, "ScriptCatalog"); ok {
		n, warns, err := explodeScripts(xmlStr, schemaRoot, sanitize)
		res.Warnings = append(res.Warnings, warns...)
		record("scripts", n, err)
	}

	// Tables: join FieldsForTables/FieldCatalog (carries fields) into <BaseTable>.
	if xmlStr, ok := loadCatalogXML(source, isDir, "FieldCatalog"); ok {
		n, err := explodeTables(xmlStr, schemaRoot)
		record("tables", n, err)
	}

	// Relationships: filename synthesized from the joined table occurrences.
	if xmlStr, ok := loadCatalogXML(source, isDir, "RelationshipCatalog"); ok {
		n, err := explodeRelationships(xmlStr, schemaRoot)
		record("relationships", n, err)
	}

	// Remaining catalogs split one object element per file (element as root).
	for _, c := range genericCatalogs {
		xmlStr, ok := loadCatalogXML(source, isDir, c.fileKey)
		if !ok {
			continue
		}
		n, err := explodeGeneric(xmlStr, c.container, c.object, c.folder, schemaRoot)
		record(c.folder, n, err)
	}

	return res, nil
}

// genericCatalog maps a FileMaker catalog to a per-object output folder for the
// straightforward "split each object element into its own file" catalogs.
type genericCatalog struct {
	fileKey   string // split-file suffix and single-file catalog element to locate
	container string // element whose direct children are the objects
	object    string // the per-object element name
	folder    string // output folder under Schema/<db>/
}

var genericCatalogs = []genericCatalog{
	{"LayoutCatalog", "LayoutCatalog", "Layout", "layouts"},
	{"TableOccurrenceCatalog", "TableOccurrenceCatalog", "TableOccurrence", "table_occurrences"},
	{"ValueListCatalog", "ValueListCatalog", "ValueList", "valuelists"},
	{"CustomFunctionsCatalog", "CustomFunctionsCatalog", "CustomFunction", "custom_functions"},
	{"CustomMenuCatalog", "CustomMenuCatalog", "CustomMenu", "custom_menus"},
	{"CustomMenuSetCatalog", "CustomMenuSetCatalog", "CustomMenuSet", "custom_menu_sets"},
	{"AccountsCatalog", "AccountsCatalog", "Account", "accounts"},
	{"PrivilegeSetsCatalog", "PrivilegeSetsCatalog", "PrivilegeSet", "privilege_sets"},
	{"ExtendedPrivilegesCatalog", "ExtendedPrivilegesCatalog", "ExtendedPrivilege", "extended_privileges"},
	{"ExternalDataSourceCatalog", "ExternalDataSourceCatalog", "ExternalDataSource", "external_data_sources"},
	{"BaseDirectoryCatalog", "BaseDirectoryCatalog", "BaseDirectory", "base_directories"},
	{"ThemeCatalog", "ThemeCatalog", "Theme", "themes"},
}

// loadCatalogXML returns the XML to parse for a catalog and whether it exists.
// A single-file export always returns the whole file (each parser locates its
// own catalog element within); a split folder returns the matching
// <…>_<fileKey>.xml, or ok=false when that catalog file is absent.
func loadCatalogXML(source string, isDir bool, fileKey string) (string, bool) {
	if !isDir {
		data, err := os.ReadFile(source)
		if err != nil {
			return "", false
		}
		return string(data), true
	}
	entries, err := os.ReadDir(source)
	if err != nil {
		return "", false
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		// FileMaker names split files "<DB>_<Catalog>.xml" (e.g.
		// Contacts_ScriptCatalog.xml); also accept a bare "<Catalog>.xml".
		if strings.HasSuffix(name, "_"+fileKey+".xml") || name == fileKey+".xml" {
			data, err := os.ReadFile(filepath.Join(source, name))
			if err != nil {
				return "", false
			}
			return string(data), true
		}
	}
	return "", false
}

// inferDatabase derives the database name from the source path: the file base
// for a single file, or the folder name for a split-catalog directory.
func inferDatabase(source string, isDir bool) string {
	base := filepath.Base(source)
	if isDir {
		return base
	}
	return strings.TrimSuffix(base, filepath.Ext(base))
}

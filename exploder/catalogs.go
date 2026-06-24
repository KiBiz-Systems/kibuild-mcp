package exploder

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

// forEachObject streams catalogXML and invokes fn once per object element that
// is a direct child of the named container element. fn receives the object's
// start-tag attributes and the object's full subtree (including its own tags)
// re-serialized as standalone XML.
//
// It tracks element depth rather than matching names, because object element
// names (Layout, TableOccurrenceReference, Script, …) also occur nested inside
// an object's own subtree; only the elements at container-depth+1 are objects.
func forEachObject(catalogXML, container, object string, fn func(attrs map[string]string, body string)) {
	dec := xml.NewDecoder(strings.NewReader(catalogXML))

	var (
		depth          int
		containerDepth = -1
		objectLevel    = -1 // depth at which objects sit; locked on first sighting
		objDepth       = -1
		buf            strings.Builder
		enc            *xml.Encoder
		attrs          map[string]string
	)

	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		switch t := tok.(type) {
		case xml.StartElement:
			depth++
			name := t.Name.Local
			switch {
			case name == container && containerDepth == -1:
				containerDepth = depth
			// Objects may be direct children of the catalog, or wrapped one level
			// deep in an <ObjectList>. Lock onto the depth of the first object seen
			// inside the catalog so both layouts are handled, while same-named
			// elements nested deeper inside an object are not mistaken for objects.
			case containerDepth != -1 && objDepth == -1 && name == object &&
				depth > containerDepth && (objectLevel == -1 || depth == objectLevel):
				if objectLevel == -1 {
					objectLevel = depth
				}
				objDepth = depth
				buf.Reset()
				enc = xml.NewEncoder(&buf)
				attrs = map[string]string{}
				for _, a := range t.Attr {
					attrs[a.Name.Local] = a.Value
				}
				_ = enc.EncodeToken(t) // object element is the root of its file
			case objDepth != -1:
				_ = enc.EncodeToken(t)
			}
		case xml.EndElement:
			name := t.Name.Local
			switch {
			case objDepth != -1 && depth == objDepth && name == object:
				_ = enc.EncodeToken(t)
				_ = enc.Flush()
				fn(attrs, buf.String())
				objDepth = -1
			case objDepth != -1:
				_ = enc.EncodeToken(t)
			case name == container && depth == containerDepth:
				containerDepth = -1
			}
			depth--
		case xml.CharData:
			if objDepth != -1 {
				_ = enc.EncodeToken(xml.CharData(append([]byte(nil), t...)))
			}
		}
	}
}

// explodeGeneric writes one file per object element (the element as its own
// root), named by its name attribute. Folder/marker entries (isFolder) are
// organizational and carry no real object, so they are skipped.
func explodeGeneric(catalogXML, container, object, folder, schemaRoot string) (int, error) {
	dir := filepath.Join(schemaRoot, folder)
	seen := map[string]int{}
	count := 0
	var writeErr error

	forEachObject(catalogXML, container, object, func(attrs map[string]string, body string) {
		if writeErr != nil {
			return
		}
		if f := attrs["isFolder"]; f == "True" || f == "Marker" {
			return
		}
		name := attrs["name"]
		if name == "" {
			name = object + "_" + attrs["id"]
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			writeErr = err
			return
		}
		fileName := uniqueName(seen, sanitizeFileName(name), attrs["id"])
		if err := os.WriteFile(filepath.Join(dir, fileName+".xml"), []byte(xmlDoc(body)), 0o644); err != nil {
			writeErr = err
			return
		}
		count++
	})
	return count, writeErr
}

// explodeTables joins FieldsForTables/FieldCatalog blocks into tables/<name>.xml.
// Each FieldCatalog block already carries <BaseTableReference name> and the
// <Field> elements that find_table parses; we just wrap it in <BaseTable>.
func explodeTables(catalogXML, schemaRoot string) (int, error) {
	dir := filepath.Join(schemaRoot, "tables")
	seen := map[string]int{}
	count := 0
	var writeErr error

	forEachObject(catalogXML, "FieldsForTables", "FieldCatalog", func(_ map[string]string, body string) {
		if writeErr != nil {
			return
		}
		name := firstAttrValue(body, "BaseTableReference", "name")
		if name == "" {
			return // a FieldCatalog with no base table reference is not a usable table
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			writeErr = err
			return
		}
		fileName := uniqueName(seen, sanitizeFileName(name), "")
		doc := xmlDoc("<BaseTable>\n" + body + "\n</BaseTable>")
		if err := os.WriteFile(filepath.Join(dir, fileName+".xml"), []byte(doc), 0o644); err != nil {
			writeErr = err
			return
		}
		count++
	})
	return count, writeErr
}

// explodeRelationships writes relationships/<name>.xml. Relationships have no
// name attribute, so the filename is synthesized from the two joined table
// occurrences (falling back to the id). inspect_relationships scans every file
// in the folder, so the filename is cosmetic.
func explodeRelationships(catalogXML, schemaRoot string) (int, error) {
	dir := filepath.Join(schemaRoot, "relationships")
	seen := map[string]int{}
	count := 0
	var writeErr error

	forEachObject(catalogXML, "RelationshipCatalog", "Relationship", func(attrs map[string]string, body string) {
		if writeErr != nil {
			return
		}
		tos := allAttrValues(body, "TableOccurrenceReference", "name")
		var name string
		if len(tos) >= 2 {
			name = tos[0] + " to " + tos[1]
		} else if id := attrs["id"]; id != "" {
			name = "Relationship " + id
		} else {
			name = "Relationship"
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			writeErr = err
			return
		}
		fileName := uniqueName(seen, sanitizeFileName(name), attrs["id"])
		if err := os.WriteFile(filepath.Join(dir, fileName+".xml"), []byte(xmlDoc(body)), 0o644); err != nil {
			writeErr = err
			return
		}
		count++
	})
	return count, writeErr
}

// xmlDoc prepends the XML declaration to an element body.
func xmlDoc(body string) string {
	return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" + body + "\n"
}

// firstAttrValue returns the value of attr on the first <element …> in xmlStr.
func firstAttrValue(xmlStr, element, attr string) string {
	vals := scanAttr(xmlStr, element, attr, true)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// allAttrValues returns attr from every <element …> in xmlStr, in order.
func allAttrValues(xmlStr, element, attr string) []string {
	return scanAttr(xmlStr, element, attr, false)
}

func scanAttr(xmlStr, element, attr string, firstOnly bool) []string {
	dec := xml.NewDecoder(strings.NewReader(xmlStr))
	var out []string
	for {
		tok, err := dec.Token()
		if err != nil {
			break
		}
		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == element {
			for _, a := range se.Attr {
				if a.Name.Local == attr {
					out = append(out, a.Value)
					if firstOnly {
						return out
					}
				}
			}
		}
	}
	return out
}

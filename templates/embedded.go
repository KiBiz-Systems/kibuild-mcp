// Package templates embeds the FileMaker XML step templates into the binary so
// the standalone server is self-contained and does not depend on the templates
// directory shipping alongside the executable.
package templates

import "embed"

//go:embed fmxml
var FS embed.FS

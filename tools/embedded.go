package tools

import "embed"

// catalogFS embeds the step catalog into the binary so the standalone server
// can validate script steps without the catalogs directory shipping alongside
// the executable. A catalog found on disk (project override) takes precedence.
//
//go:embed catalogs/step-catalog-en.json
var catalogFS embed.FS

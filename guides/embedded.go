// Package guides embeds the AI tool guide files into the binary so
// kibuild-mcp --setup can write them to the correct global config
// locations on the user's machine without requiring any extra files.
package guides

import "embed"

//go:embed *.md *.mdc
var FS embed.FS

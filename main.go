package main

import (
	"embed"

	"github.com/dillonlara115/barracuda/cmd"
)

// Embed frontend files only when building with 'serve' tag
// For API server builds (Cloud Run), we don't need embedded files
//
//go:embed web/dist
var frontendFiles embed.FS

func main() {
	// Pass embedded frontend files to cmd package
	// The 'serve' command uses these, but 'api' command ignores them
	cmd.SetFrontendFiles(frontendFiles)
	cmd.Execute()
}

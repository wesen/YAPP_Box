package docs

import (
	"embed"

	"github.com/go-go-golems/glazed/pkg/help"
)

//go:embed tutorials/*.md
var tutorialFS embed.FS

// Load registers all embedded documentation sections with the Glazed help system.
func Load(helpSystem *help.HelpSystem) error {
	return helpSystem.LoadSectionsFromFS(tutorialFS, "tutorials")
}

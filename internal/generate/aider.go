package generate

import (
	"github.com/Naly-programming/devid/internal/config"
)

// renderAider produces content for CONVENTIONS.md.
// Plain markdown, read by aider when loaded via /read or .aider.conf.yml.
func renderAider(id *config.Identity) string {
	// Aider uses plain markdown, identical format
	return renderCopilot(id)
}

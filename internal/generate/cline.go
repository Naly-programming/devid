package generate

import (
	"github.com/Naly-programming/devid/internal/config"
)

// renderCline produces content for .clinerules.
// Plain markdown, same content as the global identity.
func renderCline(id *config.Identity) string {
	// Cline uses plain markdown, identical format to copilot
	return renderCopilot(id)
}

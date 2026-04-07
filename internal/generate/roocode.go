package generate

import (
	"github.com/Naly-programming/devid/internal/config"
)

// renderRooCode produces content for .roo/rules/devid.md.
// Plain markdown, read automatically by Roo Code.
func renderRooCode(id *config.Identity) string {
	// Roo Code uses plain markdown files in .roo/rules/
	return renderCopilot(id)
}

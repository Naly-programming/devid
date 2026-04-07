package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

// renderWindsurf produces content for .windsurf/rules/devid.md.
// Markdown with YAML frontmatter for the modern Windsurf rules format.
func renderWindsurf(id *config.Identity) string {
	var b strings.Builder

	// YAML frontmatter
	b.WriteString("---\n")
	b.WriteString("trigger: always\n")
	b.WriteString("---\n\n")

	b.WriteString("# Developer Identity\n\n")

	if id.Identity.Tone != "" {
		b.WriteString(fmt.Sprintf("Tone: %s\n", id.Identity.Tone))
	}
	if id.Identity.Comments != "" {
		b.WriteString(fmt.Sprintf("Comments: %s\n", id.Identity.Comments))
	}
	if id.Identity.Responses != "" {
		b.WriteString(fmt.Sprintf("Responses: %s\n", id.Identity.Responses))
	}
	if id.Identity.Pace != "" {
		b.WriteString(fmt.Sprintf("Pace: %s\n", id.Identity.Pace))
	}

	var allStack []string
	allStack = append(allStack, id.Stack.Primary...)
	allStack = append(allStack, id.Stack.Data...)
	allStack = append(allStack, id.Stack.Infra...)
	if len(allStack) > 0 {
		b.WriteString(fmt.Sprintf("\nStack: %s\n", strings.Join(allStack, ", ")))
	}

	b.WriteString(renderAvoid(id.Stack.Avoid))
	b.WriteString(renderConventions(id.Conventions))
	b.WriteString(renderAIPrefs(id.AI))

	return b.String()
}

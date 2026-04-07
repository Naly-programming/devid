package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

// renderCursor produces an .mdc file for Cursor's .cursor/rules/ directory.
// Uses YAML frontmatter with alwaysApply: true so it loads on every session.
func renderCursor(id *config.Identity) string {
	var b strings.Builder

	// MDC frontmatter
	b.WriteString("---\n")
	b.WriteString("description: Developer identity and preferences managed by devid\n")
	b.WriteString("alwaysApply: true\n")
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

	// Stack
	var allStack []string
	allStack = append(allStack, id.Stack.Primary...)
	allStack = append(allStack, id.Stack.Data...)
	allStack = append(allStack, id.Stack.Infra...)
	if len(allStack) > 0 {
		b.WriteString(fmt.Sprintf("\nStack: %s\n", strings.Join(allStack, ", ")))
	}

	// Avoid
	if len(id.Stack.Avoid.Items) > 0 {
		b.WriteString("\nAvoid:\n")
		for _, item := range id.Stack.Avoid.Items {
			reason := id.Stack.Avoid.Reasons[strings.ReplaceAll(strings.ReplaceAll(item, "-", "_"), " ", "_")]
			if reason != "" {
				b.WriteString(fmt.Sprintf("- %s (%s)\n", item, reason))
			} else {
				b.WriteString(fmt.Sprintf("- %s\n", item))
			}
		}
	}

	// Conventions
	b.WriteString("\nConventions:\n")
	if id.Conventions.Formatting != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.Conventions.Formatting))
	}
	if id.Conventions.PRStyle != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.Conventions.PRStyle))
	}
	if id.Conventions.CommitStyle != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.Conventions.CommitStyle))
	}
	if id.Conventions.ErrorHandling != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.Conventions.ErrorHandling))
	}
	if id.Conventions.Naming != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.Conventions.Naming))
	}

	// AI preferences
	b.WriteString("\nAI:\n")
	if id.AI.Verbosity != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.AI.Verbosity))
	}
	if id.AI.Confirmation != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.AI.Confirmation))
	}
	if id.AI.Suggestions != "" {
		b.WriteString(fmt.Sprintf("- %s\n", id.AI.Suggestions))
	}
	if id.AI.Tests != "" {
		b.WriteString(fmt.Sprintf("- tests: %s\n", id.AI.Tests))
	}

	return b.String()
}

package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

func renderClaudeGlobal(id *config.Identity) string {
	var b strings.Builder

	b.WriteString("# Developer Identity\n\n")

	// Identity section
	if id.Identity.Tone != "" {
		b.WriteString(fmt.Sprintf("**Tone:** %s\n", id.Identity.Tone))
	}
	if id.Identity.Comments != "" {
		b.WriteString(fmt.Sprintf("**Comments:** %s\n", id.Identity.Comments))
	}
	if id.Identity.Responses != "" {
		b.WriteString(fmt.Sprintf("**Responses:** %s\n", id.Identity.Responses))
	}
	if id.Identity.Pace != "" {
		b.WriteString(fmt.Sprintf("**Pace:** %s\n", id.Identity.Pace))
	}

	// Stack
	var allStack []string
	allStack = append(allStack, id.Stack.Primary...)
	allStack = append(allStack, id.Stack.Data...)
	allStack = append(allStack, id.Stack.Infra...)
	if len(allStack) > 0 {
		b.WriteString(fmt.Sprintf("\n## Stack\n%s\n", joinDot(allStack)))
	}

	// Avoid
	b.WriteString(renderAvoid(id.Stack.Avoid))

	// Conventions
	b.WriteString(renderConventions(id.Conventions))

	// AI Preferences
	b.WriteString(renderAIPrefs(id.AI))

	// Learned
	if len(id.Learned.Entries) > 0 {
		b.WriteString("\n## Learned\n")
		for _, entry := range id.Learned.Entries {
			b.WriteString(fmt.Sprintf("- %s\n", entry))
		}
	}

	return b.String()
}

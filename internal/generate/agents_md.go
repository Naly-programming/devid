package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

func renderAgentsMD(id *config.Identity, proj *config.Project) string {
	var b strings.Builder

	if proj != nil {
		b.WriteString(fmt.Sprintf("# %s - Agent Instructions\n\n", proj.Name))
	} else {
		b.WriteString("# Agent Instructions\n\n")
	}

	// Identity
	if id.Identity.Tone != "" {
		b.WriteString(fmt.Sprintf("Tone: %s\n", id.Identity.Tone))
	}
	if id.Identity.Responses != "" {
		b.WriteString(fmt.Sprintf("Responses: %s\n", id.Identity.Responses))
	}
	if id.Identity.Pace != "" {
		b.WriteString(fmt.Sprintf("Pace: %s\n", id.Identity.Pace))
	}

	// Stack
	if proj != nil && len(proj.Stack) > 0 {
		b.WriteString(fmt.Sprintf("\nStack: %s\n", joinDot(proj.Stack)))
	} else {
		var allStack []string
		allStack = append(allStack, id.Stack.Primary...)
		allStack = append(allStack, id.Stack.Data...)
		if len(allStack) > 0 {
			b.WriteString(fmt.Sprintf("\nStack: %s\n", joinDot(allStack)))
		}
	}

	// Project context
	if proj != nil && proj.Context != "" {
		b.WriteString(fmt.Sprintf("\nContext: %s\n", proj.Context))
	}

	// Patterns
	if proj != nil && len(proj.Patterns) > 0 {
		b.WriteString("\nPatterns:\n")
		for _, p := range proj.Patterns {
			b.WriteString(fmt.Sprintf("- %s\n", p))
		}
	}

	// Conventions
	b.WriteString(renderConventions(id.Conventions))

	// AI preferences
	b.WriteString(renderAIPrefs(id.AI))

	return b.String()
}

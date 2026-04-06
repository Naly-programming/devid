package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

func renderClaudeProject(id *config.Identity, proj *config.Project) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", proj.Name))

	// Compressed global identity
	if id.Identity.Tone != "" {
		b.WriteString(fmt.Sprintf("**Tone:** %s\n", id.Identity.Tone))
	}
	if id.Identity.Responses != "" {
		b.WriteString(fmt.Sprintf("**Responses:** %s\n", id.Identity.Responses))
	}

	// Project context
	if proj.Context != "" {
		b.WriteString(fmt.Sprintf("\n## Context\n%s\n", proj.Context))
	}

	// Project stack
	if len(proj.Stack) > 0 {
		b.WriteString(fmt.Sprintf("\n## Stack\n%s\n", joinDot(proj.Stack)))
	}

	// Project infra
	if len(proj.Infra) > 0 {
		b.WriteString(fmt.Sprintf("\n## Infra\n%s\n", joinDot(proj.Infra)))
	}

	// Project patterns
	if len(proj.Patterns) > 0 {
		b.WriteString("\n## Patterns\n")
		for _, p := range proj.Patterns {
			b.WriteString(fmt.Sprintf("- %s\n", p))
		}
	}

	// Key conventions (compressed)
	b.WriteString(renderConventions(id.Conventions))

	// AI preferences
	b.WriteString(renderAIPrefs(id.AI))

	return b.String()
}

package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

// renderGeminiGlobal produces content for ~/.gemini/GEMINI.md (global Gemini context).
func renderGeminiGlobal(id *config.Identity) string {
	var b strings.Builder

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
		b.WriteString(fmt.Sprintf("\n## Stack\n%s\n", strings.Join(allStack, ", ")))
	}

	b.WriteString(renderAvoid(id.Stack.Avoid))
	b.WriteString(renderConventions(id.Conventions))
	b.WriteString(renderAIPrefs(id.AI))

	return b.String()
}

// renderGeminiProject produces content for {repo}/GEMINI.md (per-project).
func renderGeminiProject(id *config.Identity, proj *config.Project) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# %s\n\n", proj.Name))

	if id.Identity.Tone != "" {
		b.WriteString(fmt.Sprintf("Tone: %s\n", id.Identity.Tone))
	}

	if proj.Context != "" {
		b.WriteString(fmt.Sprintf("\n## Context\n%s\n", proj.Context))
	}
	if len(proj.Stack) > 0 {
		b.WriteString(fmt.Sprintf("\n## Stack\n%s\n", strings.Join(proj.Stack, ", ")))
	}
	if len(proj.Patterns) > 0 {
		b.WriteString("\n## Patterns\n")
		for _, p := range proj.Patterns {
			b.WriteString(fmt.Sprintf("- %s\n", p))
		}
	}

	b.WriteString(renderConventions(id.Conventions))
	b.WriteString(renderAIPrefs(id.AI))

	return b.String()
}

package generate

import (
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

func renderSnippet(id *config.Identity) string {
	var parts []string

	// Identity line
	if id.Identity.Name != "" {
		parts = append(parts, "Developer: "+id.Identity.Name+".")
	}
	if id.Identity.Tone != "" {
		parts = append(parts, "Tone: "+id.Identity.Tone+".")
	}
	if id.Identity.Responses != "" {
		parts = append(parts, "Responses: "+id.Identity.Responses+".")
	}
	if id.Identity.Pace != "" {
		parts = append(parts, "Pace: "+id.Identity.Pace+".")
	}

	// Stack
	var allStack []string
	allStack = append(allStack, id.Stack.Primary...)
	allStack = append(allStack, id.Stack.Secondary...)
	allStack = append(allStack, id.Stack.Data...)
	if len(allStack) > 0 {
		parts = append(parts, "Stack: "+strings.Join(allStack, ", ")+".")
	}

	// Avoid
	if len(id.Stack.Avoid.Items) > 0 {
		parts = append(parts, "Avoid: "+strings.Join(id.Stack.Avoid.Items, ", ")+".")
	}

	// Key conventions
	var convParts []string
	if id.Conventions.PRStyle != "" {
		convParts = append(convParts, id.Conventions.PRStyle)
	}
	if id.Conventions.CommitStyle != "" {
		convParts = append(convParts, id.Conventions.CommitStyle)
	}
	if id.Conventions.ErrorHandling != "" {
		convParts = append(convParts, id.Conventions.ErrorHandling)
	}
	if len(convParts) > 0 {
		parts = append(parts, "Conventions: "+strings.Join(convParts, "; ")+".")
	}

	// AI prefs
	var aiParts []string
	if id.AI.Verbosity != "" {
		aiParts = append(aiParts, id.AI.Verbosity)
	}
	if id.AI.Confirmation != "" {
		aiParts = append(aiParts, id.AI.Confirmation)
	}
	if id.AI.Suggestions != "" {
		aiParts = append(aiParts, id.AI.Suggestions)
	}
	if id.AI.Tests != "" {
		aiParts = append(aiParts, "tests: "+id.AI.Tests)
	}
	if len(aiParts) > 0 {
		parts = append(parts, "AI: "+strings.Join(aiParts, "; ")+".")
	}

	return strings.Join(parts, " ") + "\n"
}

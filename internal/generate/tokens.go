package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

// Token budget targets from the devid spec.
const (
	GlobalBudget  = 420  // max tokens for global context
	ProjectBudget = 80   // max additional tokens per project overlay
)

// TokenEstimate holds the estimated token count for a rendered target.
type TokenEstimate struct {
	Target string
	Tokens int
	Budget int
	Over   bool
}

// EstimateTokens returns approximate token count for a string.
// Uses chars/4 as a rough approximation for English text.
func EstimateTokens(text string) int {
	return len(text) / 4
}

// EstimateAll returns token estimates for all rendered targets.
func EstimateAll(id *config.Identity) []TokenEstimate {
	var estimates []TokenEstimate

	// Global CLAUDE.md
	global, _ := Render(id, TargetClaudeGlobal, nil)
	globalTokens := EstimateTokens(global)
	estimates = append(estimates, TokenEstimate{
		Target: "global",
		Tokens: globalTokens,
		Budget: GlobalBudget,
		Over:   globalTokens > GlobalBudget,
	})

	// Per-project overlays
	for i := range id.Projects {
		proj, _ := Render(id, TargetClaudeProject, &id.Projects[i])
		projTokens := EstimateTokens(proj)
		estimates = append(estimates, TokenEstimate{
			Target: fmt.Sprintf("project:%s", id.Projects[i].Name),
			Tokens: projTokens,
			Budget: GlobalBudget + ProjectBudget,
			Over:   projTokens > GlobalBudget+ProjectBudget,
		})
	}

	// Snippet
	snippet, _ := Render(id, TargetSnippet, nil)
	snippetTokens := EstimateTokens(snippet)
	estimates = append(estimates, TokenEstimate{
		Target: "snippet",
		Tokens: snippetTokens,
		Budget: GlobalBudget,
		Over:   snippetTokens > GlobalBudget,
	})

	return estimates
}

// FormatEstimates returns a human-readable summary of token estimates.
func FormatEstimates(estimates []TokenEstimate) string {
	var b strings.Builder
	b.WriteString("\nToken estimates:\n")
	for _, e := range estimates {
		status := "ok"
		if e.Over {
			status = "OVER BUDGET"
		}
		b.WriteString(fmt.Sprintf("  %-20s ~%d tokens (budget: %d) %s\n", e.Target, e.Tokens, e.Budget, status))
	}
	return b.String()
}

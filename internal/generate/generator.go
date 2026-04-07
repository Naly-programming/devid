package generate

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
)

// Target identifies a distribution target format.
type Target int

const (
	TargetClaudeGlobal  Target = iota // ~/.claude/CLAUDE.md
	TargetClaudeProject               // {repo}/CLAUDE.md
	TargetAgentsMD                    // {repo}/AGENTS.md
	TargetCursor                      // .cursor/rules/devid.mdc
	TargetSnippet                     // clipboard snippet
	TargetCopilot                     // .github/copilot-instructions.md
	TargetCline                       // .clinerules
	TargetRooCode                     // .roo/rules/devid.md
	TargetWindsurf                    // .windsurf/rules/devid.md
	TargetAider                       // CONVENTIONS.md
	TargetGeminiGlobal                // ~/.gemini/GEMINI.md
	TargetGeminiProject               // {repo}/GEMINI.md
)

const (
	MarkerStart = "<!-- devid:start -->"
	MarkerEnd   = "<!-- devid:end -->"
	MarkerNote  = "<!-- managed by devid - do not edit between markers -->"
)

// Render produces output for the given target.
// For project-scoped targets, pass the matching project. For global targets, project can be nil.
func Render(id *config.Identity, target Target, project *config.Project) (string, error) {
	clean := id.WithoutPrivate()
	switch target {
	case TargetClaudeGlobal:
		return renderClaudeGlobal(&clean), nil
	case TargetClaudeProject:
		if project == nil {
			return "", fmt.Errorf("project required for TargetClaudeProject")
		}
		return renderClaudeProject(&clean, project), nil
	case TargetAgentsMD:
		return renderAgentsMD(&clean, project), nil
	case TargetCursor:
		return renderCursor(&clean), nil
	case TargetSnippet:
		return renderSnippet(&clean), nil
	case TargetCopilot:
		return renderCopilot(&clean), nil
	case TargetCline:
		return renderCline(&clean), nil
	case TargetRooCode:
		return renderRooCode(&clean), nil
	case TargetWindsurf:
		return renderWindsurf(&clean), nil
	case TargetAider:
		return renderAider(&clean), nil
	case TargetGeminiGlobal:
		return renderGeminiGlobal(&clean), nil
	case TargetGeminiProject:
		if project == nil {
			return "", fmt.Errorf("project required for TargetGeminiProject")
		}
		return renderGeminiProject(&clean, project), nil
	default:
		return "", fmt.Errorf("unknown target: %d", target)
	}
}

// WrapWithMarkers wraps content between devid section markers.
func WrapWithMarkers(content string) string {
	var b strings.Builder
	b.WriteString(MarkerStart)
	b.WriteByte('\n')
	b.WriteString(MarkerNote)
	b.WriteByte('\n')
	b.WriteString(content)
	if !strings.HasSuffix(content, "\n") {
		b.WriteByte('\n')
	}
	b.WriteString(MarkerEnd)
	b.WriteByte('\n')
	return b.String()
}

// joinDot joins strings with " . " separator, skipping empty slices.
func joinDot(items []string) string {
	return strings.Join(items, " . ")
}

// renderAvoid renders the avoid section as a bullet list.
func renderAvoid(avoid config.StackAvoid) string {
	if len(avoid.Items) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\n## Avoid\n")
	for _, item := range avoid.Items {
		reason := avoid.Reasons[strings.ReplaceAll(strings.ReplaceAll(item, "-", "_"), " ", "_")]
		if reason != "" {
			b.WriteString(fmt.Sprintf("- %s (%s)\n", item, reason))
		} else {
			b.WriteString(fmt.Sprintf("- %s\n", item))
		}
	}
	return b.String()
}

// renderConventions renders conventions as a bullet list.
func renderConventions(c config.Conventions) string {
	var b strings.Builder
	b.WriteString("\n## Conventions\n")
	if c.Formatting != "" {
		b.WriteString(fmt.Sprintf("- %s\n", c.Formatting))
	}
	if c.PRStyle != "" {
		b.WriteString(fmt.Sprintf("- %s\n", c.PRStyle))
	}
	if c.CommitStyle != "" {
		b.WriteString(fmt.Sprintf("- %s\n", c.CommitStyle))
	}
	if c.ErrorHandling != "" {
		b.WriteString(fmt.Sprintf("- %s\n", c.ErrorHandling))
	}
	if c.Naming != "" {
		b.WriteString(fmt.Sprintf("- %s\n", c.Naming))
	}
	return b.String()
}

// renderAIPrefs renders AI preferences as a bullet list.
func renderAIPrefs(ai config.AI) string {
	var b strings.Builder
	b.WriteString("\n## AI Preferences\n")
	if ai.Verbosity != "" {
		b.WriteString(fmt.Sprintf("- %s\n", ai.Verbosity))
	}
	if ai.Confirmation != "" {
		b.WriteString(fmt.Sprintf("- %s\n", ai.Confirmation))
	}
	if ai.Suggestions != "" {
		b.WriteString(fmt.Sprintf("- %s\n", ai.Suggestions))
	}
	if ai.CodeComments != "" {
		b.WriteString(fmt.Sprintf("- code comments: %s\n", ai.CodeComments))
	}
	if ai.Tests != "" {
		b.WriteString(fmt.Sprintf("- tests: %s\n", ai.Tests))
	}
	return b.String()
}

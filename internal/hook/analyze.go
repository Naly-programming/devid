package hook

import (
	"fmt"
	"strings"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/extract"
)

// Signal keywords that indicate corrections, preferences, or identity-relevant statements.
var signalKeywords = []string{
	"don't", "dont", "do not",
	"stop", "never", "always",
	"prefer", "i like", "i want", "i need",
	"instead", "rather",
	"not like that", "wrong", "no,", "no.",
	"too verbose", "too formal", "too casual",
	"be more", "be less",
	"from now on", "going forward",
	"my preference", "my style",
	"i usually", "i typically",
}

const analysisSystem = `You are analyzing a coding session to identify developer preferences and style corrections.
You will receive the developer's current identity profile and a filtered set of messages from the session.
Extract ONLY novel preferences not already captured in the identity.

Rules:
- Output valid TOML only, no preamble
- Only include fields that are NEW or CHANGED
- Values must be fragments, not sentences
- If nothing novel was found, respond with exactly: NO_CHANGES
- Do not repeat things already in the current identity
- Focus on: tone corrections, workflow preferences, tool preferences, convention changes
- Ignore: project-specific instructions, one-off debugging, task-specific context`

// AnalyzeSession checks a session transcript for identity-relevant signals.
// Returns the proposed identity changes, or nil if no changes detected.
func AnalyzeSession(messages []Message, current *config.Identity) (*config.Identity, string, error) {
	if !api.Available() {
		return nil, "", fmt.Errorf("ANTHROPIC_API_KEY not set - cannot analyze session automatically")
	}

	// Filter for messages containing signal keywords
	filtered := filterForSignals(messages)
	if len(filtered) == 0 {
		return nil, "", nil // No signals, no API call
	}

	// Build the prompt
	prompt := buildAnalysisPrompt(filtered, current)

	// Call the API
	response, err := api.Call(analysisSystem, prompt)
	if err != nil {
		return nil, "", fmt.Errorf("API call failed: %w", err)
	}

	response = strings.TrimSpace(response)
	if response == "NO_CHANGES" {
		return nil, "", nil
	}

	// Parse the response
	proposed, err := extract.ParseTOMLResponse(response)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse API response: %w", err)
	}

	return proposed, response, nil
}

// filterForSignals returns messages that contain correction/preference keywords,
// along with surrounding context (the message before and after each match).
func filterForSignals(messages []Message) []Message {
	if len(messages) == 0 {
		return nil
	}

	matched := make(map[int]bool)
	for i, msg := range messages {
		if msg.Role != "user" {
			continue
		}
		lower := strings.ToLower(msg.Text)
		for _, kw := range signalKeywords {
			if strings.Contains(lower, kw) {
				// Include context: previous message, this message, next message
				if i > 0 {
					matched[i-1] = true
				}
				matched[i] = true
				if i < len(messages)-1 {
					matched[i+1] = true
				}
				break
			}
		}
	}

	if len(matched) == 0 {
		return nil
	}

	var filtered []Message
	for i := 0; i < len(messages); i++ {
		if matched[i] {
			filtered = append(filtered, messages[i])
		}
	}

	// Cap at 40 messages to keep the prompt small
	if len(filtered) > 40 {
		filtered = filtered[len(filtered)-40:]
	}

	return filtered
}

func buildAnalysisPrompt(messages []Message, current *config.Identity) string {
	var b strings.Builder

	b.WriteString("Current developer identity:\n\n```toml\n")
	if current != nil {
		tomlStr, err := extract.FormatIdentityTOML(current)
		if err == nil {
			b.WriteString(tomlStr)
		}
	}
	b.WriteString("```\n\n")

	b.WriteString("Session messages (filtered for potential preference signals):\n\n")
	for _, msg := range messages {
		b.WriteString(fmt.Sprintf("[%s]: %s\n\n", msg.Role, msg.Text))
	}

	return b.String()
}

// CountSignals returns the number of messages with signal keywords.
// Useful for deciding whether to bother with analysis.
func CountSignals(messages []Message) int {
	count := 0
	for _, msg := range messages {
		if msg.Role != "user" {
			continue
		}
		lower := strings.ToLower(msg.Text)
		for _, kw := range signalKeywords {
			if strings.Contains(lower, kw) {
				count++
				break
			}
		}
	}
	return count
}

package extract

import (
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
)

// ParseTOMLResponse extracts a TOML block from freeform text and decodes it.
// It looks for content between triple backtick fences, or assumes the whole
// input is TOML if no fences are found.
func ParseTOMLResponse(input string) (*config.Identity, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	tomlStr := extractTOMLBlock(input)

	var id config.Identity
	if _, err := toml.Decode(tomlStr, &id); err != nil {
		return nil, fmt.Errorf("failed to parse TOML: %w", err)
	}

	return &id, nil
}

// extractTOMLBlock finds TOML content between code fences, or returns the full input.
func extractTOMLBlock(input string) string {
	// Try to find fenced code block
	markers := []string{"```toml", "```"}
	for _, marker := range markers {
		start := strings.Index(input, marker)
		if start < 0 {
			continue
		}
		content := input[start+len(marker):]
		end := strings.Index(content, "```")
		if end < 0 {
			continue
		}
		return strings.TrimSpace(content[:end])
	}
	return input
}

// MergeIdentities merges an overlay onto a base identity.
// Non-zero overlay values take precedence.
func MergeIdentities(base, overlay *config.Identity) *config.Identity {
	merged := *base

	// Identity section
	if overlay.Identity.Name != "" {
		merged.Identity.Name = overlay.Identity.Name
	}
	if overlay.Identity.Tone != "" {
		merged.Identity.Tone = overlay.Identity.Tone
	}
	if overlay.Identity.Comments != "" {
		merged.Identity.Comments = overlay.Identity.Comments
	}
	if overlay.Identity.Responses != "" {
		merged.Identity.Responses = overlay.Identity.Responses
	}
	if overlay.Identity.Pace != "" {
		merged.Identity.Pace = overlay.Identity.Pace
	}

	// Stack
	if len(overlay.Stack.Primary) > 0 {
		merged.Stack.Primary = overlay.Stack.Primary
	}
	if len(overlay.Stack.Secondary) > 0 {
		merged.Stack.Secondary = overlay.Stack.Secondary
	}
	if len(overlay.Stack.Data) > 0 {
		merged.Stack.Data = overlay.Stack.Data
	}
	if len(overlay.Stack.Infra) > 0 {
		merged.Stack.Infra = overlay.Stack.Infra
	}
	if len(overlay.Stack.Messaging) > 0 {
		merged.Stack.Messaging = overlay.Stack.Messaging
	}
	if len(overlay.Stack.Cloud) > 0 {
		merged.Stack.Cloud = overlay.Stack.Cloud
	}
	if len(overlay.Stack.Testing) > 0 {
		merged.Stack.Testing = overlay.Stack.Testing
	}
	if len(overlay.Stack.Avoid.Items) > 0 {
		merged.Stack.Avoid.Items = overlay.Stack.Avoid.Items
	}
	if len(overlay.Stack.Avoid.Reasons) > 0 {
		if merged.Stack.Avoid.Reasons == nil {
			merged.Stack.Avoid.Reasons = make(map[string]string)
		}
		for k, v := range overlay.Stack.Avoid.Reasons {
			merged.Stack.Avoid.Reasons[k] = v
		}
	}

	// Conventions
	if overlay.Conventions.Formatting != "" {
		merged.Conventions.Formatting = overlay.Conventions.Formatting
	}
	if overlay.Conventions.PRStyle != "" {
		merged.Conventions.PRStyle = overlay.Conventions.PRStyle
	}
	if overlay.Conventions.CommitStyle != "" {
		merged.Conventions.CommitStyle = overlay.Conventions.CommitStyle
	}
	if overlay.Conventions.ErrorHandling != "" {
		merged.Conventions.ErrorHandling = overlay.Conventions.ErrorHandling
	}
	if overlay.Conventions.Naming != "" {
		merged.Conventions.Naming = overlay.Conventions.Naming
	}

	// AI
	if overlay.AI.Verbosity != "" {
		merged.AI.Verbosity = overlay.AI.Verbosity
	}
	if overlay.AI.Confirmation != "" {
		merged.AI.Confirmation = overlay.AI.Confirmation
	}
	if overlay.AI.Suggestions != "" {
		merged.AI.Suggestions = overlay.AI.Suggestions
	}
	if overlay.AI.CodeComments != "" {
		merged.AI.CodeComments = overlay.AI.CodeComments
	}
	if overlay.AI.Tests != "" {
		merged.AI.Tests = overlay.AI.Tests
	}

	// Learned - append new entries
	if len(overlay.Learned.Entries) > 0 {
		existing := make(map[string]bool)
		for _, e := range merged.Learned.Entries {
			existing[e] = true
		}
		for _, e := range overlay.Learned.Entries {
			if !existing[e] {
				merged.Learned.Entries = append(merged.Learned.Entries, e)
			}
		}
	}

	return &merged
}

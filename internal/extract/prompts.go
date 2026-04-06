package extract

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
)

// ExtractionPrompt is the canonical prompt for extracting developer identity.
const ExtractionPrompt = `You are helping build a compressed developer identity file.
Answer each section in fragments not sentences. Be terse.
No explanations unless asked. Output valid TOML only.

Extract or infer the following about this developer:

[identity]
- name
- tone (how they communicate, 5 words max)
- comment style preference
- response format preference

[stack]
- primary languages/frameworks
- secondary languages/frameworks
- data/infra tools
- things they explicitly avoid and why

[conventions]
- formatting preferences
- PR style
- commit style
- error handling approach
- naming conventions

[ai]
- verbosity preference
- how much hand-holding they want
- whether to challenge or comply
- test writing preference

[learned]
- any explicit preferences stated in this session not covered above

Rules:
- Values must be fragments, not sentences
- Lists preferred over prose
- If unsure, omit rather than guess
- Output the TOML block only, no preamble`

// BuildSyncPrompt builds a contextual prompt that includes the current TOML.
func BuildSyncPrompt(current *config.Identity) string {
	var buf bytes.Buffer
	buf.WriteString(ExtractionPrompt)
	buf.WriteString("\n\n---\n\nCurrent identity.toml for reference (update or add to it, do not repeat unchanged values):\n\n```toml\n")
	if current != nil {
		toml.NewEncoder(&buf).Encode(current)
	}
	buf.WriteString("```\n")
	return buf.String()
}

// FormatIdentityTOML returns the TOML representation of an identity.
func FormatIdentityTOML(id *config.Identity) (string, error) {
	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(id); err != nil {
		return "", fmt.Errorf("failed to encode identity: %w", err)
	}
	return buf.String(), nil
}

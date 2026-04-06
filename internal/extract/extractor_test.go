package extract

import (
	"testing"

	"github.com/Naly-programming/devid/internal/config"
)

func TestParseTOMLResponseFenced(t *testing.T) {
	input := "Here's your identity:\n\n```toml\n[identity]\nname = \"Test\"\ntone = \"direct\"\n```\n\nLet me know if you need changes."

	id, err := ParseTOMLResponse(input)
	if err != nil {
		t.Fatalf("ParseTOMLResponse failed: %v", err)
	}
	if id.Identity.Name != "Test" {
		t.Errorf("name = %q, want %q", id.Identity.Name, "Test")
	}
	if id.Identity.Tone != "direct" {
		t.Errorf("tone = %q, want %q", id.Identity.Tone, "direct")
	}
}

func TestParseTOMLResponseRaw(t *testing.T) {
	input := `[identity]
name = "Test"

[stack]
primary = ["Go", "TypeScript"]
`

	id, err := ParseTOMLResponse(input)
	if err != nil {
		t.Fatalf("ParseTOMLResponse failed: %v", err)
	}
	if id.Identity.Name != "Test" {
		t.Errorf("name = %q, want %q", id.Identity.Name, "Test")
	}
	if len(id.Stack.Primary) != 2 {
		t.Errorf("stack.primary length = %d, want 2", len(id.Stack.Primary))
	}
}

func TestParseTOMLResponseEmpty(t *testing.T) {
	_, err := ParseTOMLResponse("")
	if err == nil {
		t.Error("expected error on empty input")
	}
}

func TestParseTOMLResponseGarbage(t *testing.T) {
	_, err := ParseTOMLResponse("this is not TOML at all = = = [[[")
	if err == nil {
		t.Error("expected error on garbage input")
	}
}

func TestMergeIdentities(t *testing.T) {
	base := &config.Identity{
		Identity: config.IdentitySection{
			Name: "Nathan",
			Tone: "direct",
			Pace: "fast",
		},
		Stack: config.Stack{
			Primary: []string{"Go", "TypeScript"},
		},
		AI: config.AI{
			Verbosity: "concise",
			Tests:     "write them",
		},
		Learned: config.Learned{
			Entries: []string{"prefers explicit errors"},
		},
	}

	overlay := &config.Identity{
		Identity: config.IdentitySection{
			Tone: "very direct", // Changed
			// Name and Pace empty - should keep base values
		},
		Stack: config.Stack{
			Primary: []string{"Go", "TypeScript", "Rust"}, // Updated
		},
		Learned: config.Learned{
			Entries: []string{
				"prefers explicit errors", // Duplicate - should not appear twice
				"likes table tests",       // New
			},
		},
	}

	merged := MergeIdentities(base, overlay)

	if merged.Identity.Name != "Nathan" {
		t.Errorf("name should be preserved: got %q", merged.Identity.Name)
	}
	if merged.Identity.Tone != "very direct" {
		t.Errorf("tone should be overridden: got %q", merged.Identity.Tone)
	}
	if merged.Identity.Pace != "fast" {
		t.Errorf("pace should be preserved: got %q", merged.Identity.Pace)
	}
	if len(merged.Stack.Primary) != 3 {
		t.Errorf("stack.primary should be updated: got %v", merged.Stack.Primary)
	}
	if merged.AI.Verbosity != "concise" {
		t.Errorf("ai.verbosity should be preserved: got %q", merged.AI.Verbosity)
	}
	if len(merged.Learned.Entries) != 2 {
		t.Errorf("learned should have 2 entries (no dups): got %v", merged.Learned.Entries)
	}
}

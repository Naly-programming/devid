package config

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"
)

func loadExample(t *testing.T) Identity {
	t.Helper()
	var id Identity
	_, err := toml.DecodeFile("../../schema/identity.toml.example", &id)
	if err != nil {
		t.Fatalf("failed to decode example: %v", err)
	}
	return id
}

func TestDecodeExample(t *testing.T) {
	id := loadExample(t)

	if id.Meta.Version != "1" {
		t.Errorf("meta.version = %q, want %q", id.Meta.Version, "1")
	}
	if id.Identity.Name != "Nathan" {
		t.Errorf("identity.name = %q, want %q", id.Identity.Name, "Nathan")
	}
	if len(id.Stack.Primary) != 3 {
		t.Errorf("stack.primary length = %d, want 3", len(id.Stack.Primary))
	}
	if id.Stack.Primary[0] != "Go" {
		t.Errorf("stack.primary[0] = %q, want %q", id.Stack.Primary[0], "Go")
	}
	if len(id.Stack.Avoid.Items) != 3 {
		t.Errorf("stack.avoid.items length = %d, want 3", len(id.Stack.Avoid.Items))
	}
	if id.Stack.Avoid.Reasons["playwright_go"] != "community wrapper, lags behind, use JS/TS for Playwright" {
		t.Errorf("stack.avoid.reasons.playwright_go = %v", id.Stack.Avoid.Reasons["playwright_go"])
	}
	if len(id.Projects) != 2 {
		t.Errorf("projects length = %d, want 2", len(id.Projects))
	}
	if id.Projects[0].Name != "coentry" {
		t.Errorf("projects[0].name = %q, want %q", id.Projects[0].Name, "coentry")
	}
	if id.Conventions.PRStyle != "small focused PRs, one concern per PR" {
		t.Errorf("conventions.pr_style = %q", id.Conventions.PRStyle)
	}
	if id.AI.Verbosity != "concise, skip preamble, get to the point" {
		t.Errorf("ai.verbosity = %q", id.AI.Verbosity)
	}
	if len(id.Learned.Entries) != 1 {
		t.Errorf("learned.entries length = %d, want 1", len(id.Learned.Entries))
	}
	if id.Private == nil {
		t.Error("private section should not be nil")
	}
	if id.Private["api_key"] != "example-secret-that-should-never-appear-in-output" {
		t.Errorf("private.api_key = %v", id.Private["api_key"])
	}
}

func TestRoundTrip(t *testing.T) {
	original := loadExample(t)

	var buf bytes.Buffer
	if err := toml.NewEncoder(&buf).Encode(original); err != nil {
		t.Fatalf("failed to encode: %v", err)
	}

	var decoded Identity
	if _, err := toml.Decode(buf.String(), &decoded); err != nil {
		t.Fatalf("failed to re-decode: %v", err)
	}

	// Spot-check key fields survive the round trip
	if decoded.Identity.Name != original.Identity.Name {
		t.Errorf("name: got %q, want %q", decoded.Identity.Name, original.Identity.Name)
	}
	if len(decoded.Stack.Primary) != len(original.Stack.Primary) {
		t.Errorf("stack.primary length: got %d, want %d", len(decoded.Stack.Primary), len(original.Stack.Primary))
	}
	if len(decoded.Projects) != len(original.Projects) {
		t.Errorf("projects length: got %d, want %d", len(decoded.Projects), len(original.Projects))
	}
	if decoded.Conventions.CommitStyle != original.Conventions.CommitStyle {
		t.Errorf("commit_style: got %q, want %q", decoded.Conventions.CommitStyle, original.Conventions.CommitStyle)
	}
	if len(decoded.Stack.Avoid.Items) != len(original.Stack.Avoid.Items) {
		t.Errorf("avoid.items length: got %d, want %d", len(decoded.Stack.Avoid.Items), len(original.Stack.Avoid.Items))
	}
}

func TestEmptyFile(t *testing.T) {
	var id Identity
	if _, err := toml.Decode("", &id); err != nil {
		t.Fatalf("empty TOML should not error: %v", err)
	}
	if id.Identity.Name != "" {
		t.Errorf("identity.name should be empty, got %q", id.Identity.Name)
	}
	if id.Projects != nil {
		t.Errorf("projects should be nil, got %v", id.Projects)
	}
}

func TestPartialFile(t *testing.T) {
	input := `
[identity]
name = "Test"
tone = "chill"
`
	var id Identity
	if _, err := toml.Decode(input, &id); err != nil {
		t.Fatalf("partial TOML should not error: %v", err)
	}
	if id.Identity.Name != "Test" {
		t.Errorf("identity.name = %q, want %q", id.Identity.Name, "Test")
	}
	if id.Identity.Tone != "chill" {
		t.Errorf("identity.tone = %q, want %q", id.Identity.Tone, "chill")
	}
	// Other sections should be zero-value
	if len(id.Stack.Primary) != 0 {
		t.Errorf("stack.primary should be empty, got %v", id.Stack.Primary)
	}
	if id.Conventions.PRStyle != "" {
		t.Errorf("conventions.pr_style should be empty, got %q", id.Conventions.PRStyle)
	}
}

func TestWithoutPrivate(t *testing.T) {
	id := loadExample(t)
	if id.Private == nil {
		t.Fatal("expected private section in example")
	}

	clean := id.WithoutPrivate()
	if clean.Private != nil {
		t.Error("WithoutPrivate should clear the private section")
	}
	// Original should be untouched
	if id.Private == nil {
		t.Error("original identity should still have private section")
	}
}

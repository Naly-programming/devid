package sync

import (
	"testing"

	"github.com/Naly-programming/devid/internal/config"
)

func TestDiffIdentical(t *testing.T) {
	id := &config.Identity{
		Identity: config.IdentitySection{Name: "Nathan", Tone: "direct"},
	}

	diff, err := DiffIdentities(id, id)
	if err != nil {
		t.Fatalf("DiffIdentities failed: %v", err)
	}

	// All lines should be "equal" (prefixed with "  ")
	for _, line := range splitLines(diff) {
		if line != "" && line[0] != ' ' {
			t.Errorf("expected all equal lines, got: %q", line)
		}
	}
}

func TestDiffChanged(t *testing.T) {
	current := &config.Identity{
		Identity: config.IdentitySection{Name: "Nathan", Tone: "direct"},
	}
	proposed := &config.Identity{
		Identity: config.IdentitySection{Name: "Nathan", Tone: "very direct"},
	}

	diff, err := DiffIdentities(current, proposed)
	if err != nil {
		t.Fatalf("DiffIdentities failed: %v", err)
	}

	if diff == "" {
		t.Error("diff should not be empty for changed identities")
	}

	// Should contain both deletions and insertions
	hasDelete := false
	hasInsert := false
	for _, line := range splitLines(diff) {
		if len(line) > 0 && line[0] == '-' {
			hasDelete = true
		}
		if len(line) > 0 && line[0] == '+' {
			hasInsert = true
		}
	}
	if !hasDelete || !hasInsert {
		t.Errorf("diff should have both deletions and insertions, got:\n%s", diff)
	}
}

func TestDiffAddedLearned(t *testing.T) {
	current := &config.Identity{
		Learned: config.Learned{Entries: []string{"entry one"}},
	}
	proposed := &config.Identity{
		Learned: config.Learned{Entries: []string{"entry one", "entry two"}},
	}

	diff, err := DiffIdentities(current, proposed)
	if err != nil {
		t.Fatalf("DiffIdentities failed: %v", err)
	}
	if diff == "" {
		t.Error("diff should not be empty when entries are added")
	}
}

func TestApplyCandidate(t *testing.T) {
	current := &config.Identity{
		Identity: config.IdentitySection{Name: "Nathan", Tone: "direct"},
		AI:       config.AI{Verbosity: "concise"},
	}

	candidate := Candidate{
		Proposed: &config.Identity{
			Identity: config.IdentitySection{Tone: "very direct"},
		},
	}

	merged := ApplyCandidate(current, candidate)
	if merged.Identity.Name != "Nathan" {
		t.Errorf("name should be preserved: got %q", merged.Identity.Name)
	}
	if merged.Identity.Tone != "very direct" {
		t.Errorf("tone should be updated: got %q", merged.Identity.Tone)
	}
	if merged.AI.Verbosity != "concise" {
		t.Errorf("verbosity should be preserved: got %q", merged.AI.Verbosity)
	}
}

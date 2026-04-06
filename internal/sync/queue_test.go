package sync

import (
	"os"
	"testing"
	"time"

	"github.com/Naly-programming/devid/internal/config"
)

func setupTestHome(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	config.SetHomeDir(func() (string, error) { return dir, nil })
	t.Cleanup(func() { config.SetHomeDir(os.UserHomeDir) })
}

func TestEnqueueAndList(t *testing.T) {
	setupTestHome(t)

	c1 := Candidate{
		Timestamp: time.Unix(1000, 0),
		Source:    "sync",
		Proposed: &config.Identity{
			Identity: config.IdentitySection{Name: "First"},
		},
		Diff: "some diff",
	}
	c2 := Candidate{
		Timestamp: time.Unix(2000, 0),
		Source:    "manual",
		Proposed: &config.Identity{
			Identity: config.IdentitySection{Name: "Second"},
		},
	}

	if err := Enqueue(c1); err != nil {
		t.Fatalf("Enqueue c1 failed: %v", err)
	}
	if err := Enqueue(c2); err != nil {
		t.Fatalf("Enqueue c2 failed: %v", err)
	}

	candidates, err := ListQueue()
	if err != nil {
		t.Fatalf("ListQueue failed: %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(candidates))
	}

	// Should be sorted by timestamp
	if candidates[0].Proposed.Identity.Name != "First" {
		t.Errorf("first candidate name = %q, want %q", candidates[0].Proposed.Identity.Name, "First")
	}
	if candidates[1].Proposed.Identity.Name != "Second" {
		t.Errorf("second candidate name = %q, want %q", candidates[1].Proposed.Identity.Name, "Second")
	}
	if candidates[0].Source != "sync" {
		t.Errorf("first source = %q, want %q", candidates[0].Source, "sync")
	}
	if candidates[0].Diff != "some diff" {
		t.Errorf("first diff = %q, want %q", candidates[0].Diff, "some diff")
	}
}

func TestEmptyQueue(t *testing.T) {
	setupTestHome(t)

	candidates, err := ListQueue()
	if err != nil {
		t.Fatalf("ListQueue failed: %v", err)
	}
	if len(candidates) != 0 {
		t.Errorf("expected 0 candidates, got %d", len(candidates))
	}
}

func TestRemoveCandidate(t *testing.T) {
	setupTestHome(t)

	c := Candidate{
		Timestamp: time.Unix(3000, 0),
		Source:    "sync",
		Proposed:  &config.Identity{},
	}
	if err := Enqueue(c); err != nil {
		t.Fatalf("Enqueue failed: %v", err)
	}

	if err := RemoveCandidate(3000); err != nil {
		t.Fatalf("RemoveCandidate failed: %v", err)
	}

	candidates, err := ListQueue()
	if err != nil {
		t.Fatalf("ListQueue failed: %v", err)
	}
	if len(candidates) != 0 {
		t.Errorf("expected 0 candidates after remove, got %d", len(candidates))
	}
}

func TestClearQueue(t *testing.T) {
	setupTestHome(t)

	for i := int64(0); i < 3; i++ {
		c := Candidate{
			Timestamp: time.Unix(i*1000+1000, 0),
			Source:    "sync",
			Proposed:  &config.Identity{},
		}
		Enqueue(c)
	}

	if err := ClearQueue(); err != nil {
		t.Fatalf("ClearQueue failed: %v", err)
	}

	candidates, _ := ListQueue()
	if len(candidates) != 0 {
		t.Errorf("expected 0 candidates after clear, got %d", len(candidates))
	}
}

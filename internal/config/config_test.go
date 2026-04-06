package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func setupTestHome(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	SetHomeDir(func() (string, error) { return dir, nil })
	t.Cleanup(func() { SetHomeDir(os.UserHomeDir) })
	return dir
}

func TestSaveAndLoad(t *testing.T) {
	setupTestHome(t)

	original := &Identity{
		Meta: Meta{Version: "1"},
		Identity: IdentitySection{
			Name: "Test",
			Tone: "direct",
		},
		Stack: Stack{
			Primary: []string{"Go", "TypeScript"},
		},
		Conventions: Conventions{
			PRStyle: "small PRs",
		},
		AI: AI{
			Verbosity: "concise",
		},
	}

	if err := Save(original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Identity.Name != "Test" {
		t.Errorf("name = %q, want %q", loaded.Identity.Name, "Test")
	}
	if loaded.Identity.Tone != "direct" {
		t.Errorf("tone = %q, want %q", loaded.Identity.Tone, "direct")
	}
	if len(loaded.Stack.Primary) != 2 {
		t.Errorf("stack.primary length = %d, want 2", len(loaded.Stack.Primary))
	}
	if loaded.Meta.UpdatedAt.IsZero() {
		t.Error("updated_at should be set after save")
	}
}

func TestLoadMissing(t *testing.T) {
	setupTestHome(t)

	_, err := Load()
	if !errors.Is(err, ErrNoIdentity) {
		t.Errorf("Load on missing file: got %v, want ErrNoIdentity", err)
	}
}

func TestAtomicWrite(t *testing.T) {
	home := setupTestHome(t)

	id := &Identity{Meta: Meta{Version: "1"}, Identity: IdentitySection{Name: "Test"}}
	if err := Save(id); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	tmp := filepath.Join(home, ".devid", "identity.toml.tmp")
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Error("temp file should not exist after save")
	}

	p := filepath.Join(home, ".devid", "identity.toml")
	if _, err := os.Stat(p); err != nil {
		t.Errorf("identity.toml should exist: %v", err)
	}
}

func TestCreatesDirectory(t *testing.T) {
	home := setupTestHome(t)

	devidDir := filepath.Join(home, ".devid")
	if _, err := os.Stat(devidDir); !os.IsNotExist(err) {
		t.Fatal(".devid dir should not exist before save")
	}

	id := &Identity{Meta: Meta{Version: "1"}}
	if err := Save(id); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if _, err := os.Stat(devidDir); err != nil {
		t.Errorf(".devid dir should exist after save: %v", err)
	}
}

func TestExists(t *testing.T) {
	setupTestHome(t)

	if Exists() {
		t.Error("Exists should return false before init")
	}

	id := &Identity{Meta: Meta{Version: "1"}}
	if err := Save(id); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	if !Exists() {
		t.Error("Exists should return true after save")
	}
}

func TestQueueDir(t *testing.T) {
	home := setupTestHome(t)

	dir, err := QueueDir()
	if err != nil {
		t.Fatalf("QueueDir failed: %v", err)
	}

	expected := filepath.Join(home, ".devid", "queue")
	if dir != expected {
		t.Errorf("QueueDir = %q, want %q", dir, expected)
	}

	if _, err := os.Stat(dir); err != nil {
		t.Errorf("queue dir should exist: %v", err)
	}
}

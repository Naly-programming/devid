package distribute

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteFileNew(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.mdc")

	action, err := writeFile(path, "hello\n")
	if err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}
	if action != "created" {
		t.Errorf("action = %q, want %q", action, "created")
	}

	data, _ := os.ReadFile(path)
	if string(data) != "hello\n" {
		t.Errorf("content = %q", string(data))
	}
}

func TestWriteFileUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.mdc")

	os.WriteFile(path, []byte("hello\n"), 0o644)

	action, err := writeFile(path, "hello\n")
	if err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}
	if action != "unchanged" {
		t.Errorf("action = %q, want %q", action, "unchanged")
	}
}

func TestWriteFileUpdated(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.mdc")

	os.WriteFile(path, []byte("old\n"), 0o644)

	action, err := writeFile(path, "new\n")
	if err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}
	if action != "updated" {
		t.Errorf("action = %q, want %q", action, "updated")
	}

	data, _ := os.ReadFile(path)
	if string(data) != "new\n" {
		t.Errorf("content = %q", string(data))
	}
}

package distribute

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
)

func testIdentity() *config.Identity {
	return &config.Identity{
		Meta: config.Meta{Version: "1"},
		Identity: config.IdentitySection{
			Name:      "Test",
			Tone:      "direct",
			Comments:  "terse",
			Responses: "prose",
			Pace:      "fast",
		},
		Stack: config.Stack{
			Primary: []string{"Go", "TypeScript"},
			Avoid: config.StackAvoid{
				Items: []string{"ORM"},
			},
		},
		Conventions: config.Conventions{
			PRStyle:     "small PRs",
			CommitStyle: "conventional",
		},
		AI: config.AI{
			Verbosity: "concise",
			Tests:     "write them",
		},
		Projects: []config.Project{
			{
				Name:    "myproject",
				Repo:    "myproject",
				Stack:   []string{"Go"},
				Context: "test project",
			},
		},
	}
}

func TestWriteWithMarkersNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	action, err := writeWithMarkers(path, "hello world\n")
	if err != nil {
		t.Fatalf("writeWithMarkers failed: %v", err)
	}
	if action != "created" {
		t.Errorf("action = %q, want %q", action, "created")
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, generate.MarkerStart) {
		t.Error("should contain marker start")
	}
	if !strings.Contains(content, generate.MarkerEnd) {
		t.Error("should contain marker end")
	}
	if !strings.Contains(content, "hello world") {
		t.Error("should contain the content")
	}
}

func TestWriteWithMarkersExistingWithMarkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	// Create file with existing markers and surrounding content
	initial := "# My Notes\n\n" +
		generate.MarkerStart + "\n" +
		"old content\n" +
		generate.MarkerEnd + "\n" +
		"\n## My Custom Section\nkeep this\n"

	os.WriteFile(path, []byte(initial), 0o644)

	action, err := writeWithMarkers(path, "new content\n")
	if err != nil {
		t.Fatalf("writeWithMarkers failed: %v", err)
	}
	if action != "updated" {
		t.Errorf("action = %q, want %q", action, "updated")
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	if !strings.Contains(content, "# My Notes") {
		t.Error("should preserve content before markers")
	}
	if !strings.Contains(content, "new content") {
		t.Error("should contain new content")
	}
	if strings.Contains(content, "old content") {
		t.Error("should not contain old content")
	}
	if !strings.Contains(content, "My Custom Section") {
		t.Error("should preserve content after markers")
	}
	if !strings.Contains(content, "keep this") {
		t.Error("should preserve content after markers")
	}
}

func TestWriteWithMarkersExistingNoMarkers(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	original := "# Existing Content\nsome stuff here\n"
	os.WriteFile(path, []byte(original), 0o644)

	action, err := writeWithMarkers(path, "devid content\n")
	if err != nil {
		t.Fatalf("writeWithMarkers failed: %v", err)
	}
	if action != "updated" {
		t.Errorf("action = %q, want %q", action, "updated")
	}

	data, _ := os.ReadFile(path)
	content := string(data)

	// Markers should be prepended, original content after
	markerIdx := strings.Index(content, generate.MarkerStart)
	originalIdx := strings.Index(content, "# Existing Content")
	if markerIdx < 0 || originalIdx < 0 {
		t.Fatal("both marker and original content should exist")
	}
	if markerIdx >= originalIdx {
		t.Error("marker block should come before original content")
	}
	if !strings.Contains(content, "some stuff here") {
		t.Error("original content should be preserved")
	}
}

func TestWriteWithMarkersUnchanged(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.md")

	content := "test content\n"
	// Write once
	writeWithMarkers(path, content)
	// Write again with same content
	action, err := writeWithMarkers(path, content)
	if err != nil {
		t.Fatalf("writeWithMarkers failed: %v", err)
	}
	if action != "unchanged" {
		t.Errorf("action = %q, want %q", action, "unchanged")
	}
}

func TestMatchProject(t *testing.T) {
	id := testIdentity()

	proj := matchProject(id, "myproject")
	if proj == nil {
		t.Fatal("should match myproject")
	}
	if proj.Name != "myproject" {
		t.Errorf("name = %q, want %q", proj.Name, "myproject")
	}

	proj = matchProject(id, "nonexistent")
	if proj != nil {
		t.Error("should not match nonexistent")
	}

	// Case insensitive
	proj = matchProject(id, "MyProject")
	if proj == nil {
		t.Error("should match case-insensitively")
	}
}

func TestDistributeCollectsResults(t *testing.T) {
	dir := t.TempDir()

	// Override home dir for global target
	config.SetHomeDir(func() (string, error) { return dir, nil })
	t.Cleanup(func() { config.SetHomeDir(os.UserHomeDir) })

	// Override repo detector to point to temp dir
	repoRoot := filepath.Join(dir, "repo")
	os.MkdirAll(repoRoot, 0o755)
	SetRepoDetector(func() (string, string, error) {
		return repoRoot, "myproject", nil
	})
	t.Cleanup(func() { SetRepoDetector(detectRepo) })

	id := testIdentity()
	results := Distribute(id)

	// Should have: global, project, agents, cursor
	if len(results) < 4 {
		t.Fatalf("expected at least 4 results, got %d", len(results))
	}

	targets := make(map[string]bool)
	for _, r := range results {
		targets[r.Target] = true
		if r.Err != nil {
			t.Errorf("target %q failed: %v", r.Target, r.Err)
		}
		if r.Action == "" {
			t.Errorf("target %q has empty action", r.Target)
		}
	}

	for _, expected := range []string{"claude-global", "claude-project", "agents-md", "cursor"} {
		if !targets[expected] {
			t.Errorf("missing target: %s", expected)
		}
	}
}

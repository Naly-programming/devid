package scan

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFindExistingContextFiles(t *testing.T) {
	dir := t.TempDir()

	// Create a fake repo with a CLAUDE.md
	repoDir := filepath.Join(dir, "myrepo")
	os.MkdirAll(repoDir, 0o755)
	os.WriteFile(filepath.Join(repoDir, "CLAUDE.md"), []byte("# My Project\nUse Go.\n"), 0o644)

	sources := FindExistingContextFiles([]string{dir})

	found := false
	for _, src := range sources {
		if strings.Contains(src.Path, "myrepo") {
			found = true
			if !strings.Contains(src.Content, "Use Go") {
				t.Errorf("expected content to contain 'Use Go', got %q", src.Content)
			}
		}
	}
	if !found {
		t.Error("expected to find myrepo CLAUDE.md")
	}
}

func TestStripDevidMarkers(t *testing.T) {
	input := "# My Notes\n\n<!-- devid:start -->\ngenerated stuff\n<!-- devid:end -->\n\n## Custom\nkeep this\n"
	result := stripDevidMarkers(input)

	if strings.Contains(result, "generated stuff") {
		t.Error("should strip content between markers")
	}
	if !strings.Contains(result, "My Notes") {
		t.Error("should keep content before markers")
	}
	if !strings.Contains(result, "keep this") {
		t.Error("should keep content after markers")
	}
}

func TestStripDevidMarkersNoMarkers(t *testing.T) {
	input := "# Just a normal file\nNo markers here.\n"
	result := stripDevidMarkers(input)
	if result != input {
		t.Errorf("should return input unchanged, got %q", result)
	}
}

func TestBuildInferencePrompt(t *testing.T) {
	sources := []InferredSource{
		{Path: "/home/user/repo1/CLAUDE.md", Content: "Use Go. Be concise."},
		{Path: "/home/user/repo2/CLAUDE.md", Content: "TypeScript project. No ORMs."},
	}

	prompt := BuildInferencePrompt(sources)

	if !strings.Contains(prompt, "repo1/CLAUDE.md") {
		t.Error("should contain file path")
	}
	if !strings.Contains(prompt, "Use Go") {
		t.Error("should contain file content")
	}
	if !strings.Contains(prompt, "No ORMs") {
		t.Error("should contain second file content")
	}
}

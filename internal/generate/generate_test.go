package generate

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
)

var update = flag.Bool("update", false, "update golden files")

func loadTestIdentity(t *testing.T) *config.Identity {
	t.Helper()
	var id config.Identity
	_, err := toml.DecodeFile("../../schema/identity.toml.example", &id)
	if err != nil {
		t.Fatalf("failed to decode example: %v", err)
	}
	return &id
}

func goldenPath(name string) string {
	return filepath.Join("..", "..", "testdata", "golden", name)
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := goldenPath(name)

	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("failed to create golden dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("failed to write golden file: %v", err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden file %s not found (run with -update to create): %v", name, err)
	}
	if got != string(want) {
		t.Errorf("output does not match golden file %s:\n--- got ---\n%s\n--- want ---\n%s", name, got, string(want))
	}
}

func TestRenderClaudeGlobal(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetClaudeGlobal, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "claude_global.md", out)
}

func TestRenderClaudeProject(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetClaudeProject, &id.Projects[0])
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "claude_project.md", out)
}

func TestRenderClaudeProjectNilError(t *testing.T) {
	id := loadTestIdentity(t)
	_, err := Render(id, TargetClaudeProject, nil)
	if err == nil {
		t.Error("expected error when project is nil")
	}
}

func TestRenderAgentsMD(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetAgentsMD, &id.Projects[0])
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "agents.md", out)
}

func TestRenderAgentsMDGlobal(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetAgentsMD, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "agents_global.md", out)
}

func TestRenderCursor(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetCursor, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "cursor.txt", out)
}

func TestRenderSnippet(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetSnippet, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	assertGolden(t, "snippet.txt", out)
}

func TestPrivateDataExcluded(t *testing.T) {
	id := loadTestIdentity(t)
	if id.Private == nil {
		t.Fatal("example should have private section")
	}

	targets := []struct {
		name    Target
		project *config.Project
	}{
		{TargetClaudeGlobal, nil},
		{TargetClaudeProject, &id.Projects[0]},
		{TargetAgentsMD, &id.Projects[0]},
		{TargetCursor, nil},
		{TargetSnippet, nil},
	}

	for _, tc := range targets {
		out, err := Render(id, tc.name, tc.project)
		if err != nil {
			t.Errorf("target %d: render failed: %v", tc.name, err)
			continue
		}
		if strings.Contains(out, "example-secret") {
			t.Errorf("target %d: private data leaked into output", tc.name)
		}
		if strings.Contains(out, "api_key") {
			t.Errorf("target %d: private key name leaked into output", tc.name)
		}
	}
}

func TestWrapWithMarkers(t *testing.T) {
	content := "hello\nworld\n"
	wrapped := WrapWithMarkers(content)

	if !strings.HasPrefix(wrapped, MarkerStart+"\n") {
		t.Error("should start with MarkerStart")
	}
	if !strings.Contains(wrapped, MarkerNote) {
		t.Error("should contain marker note")
	}
	if !strings.Contains(wrapped, "hello\nworld\n") {
		t.Error("should contain original content")
	}
	if !strings.HasSuffix(wrapped, MarkerEnd+"\n") {
		t.Error("should end with MarkerEnd")
	}
}

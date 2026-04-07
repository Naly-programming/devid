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

func TestRenderCopilot(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetCopilot, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(out, "Developer Identity") {
		t.Error("expected Developer Identity heading")
	}
	if !strings.Contains(out, "Go") {
		t.Error("expected Go in stack")
	}
}

func TestRenderCline(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetCline, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if out == "" {
		t.Error("empty output")
	}
}

func TestRenderRooCode(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetRooCode, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if out == "" {
		t.Error("empty output")
	}
}

func TestRenderWindsurf(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetWindsurf, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.HasPrefix(out, "---\n") {
		t.Error("windsurf should have YAML frontmatter")
	}
	if !strings.Contains(out, "trigger: always") {
		t.Error("windsurf should have trigger: always")
	}
}

func TestRenderAider(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetAider, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if out == "" {
		t.Error("empty output")
	}
}

func TestRenderGeminiGlobal(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetGeminiGlobal, nil)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(out, "Developer Identity") {
		t.Error("expected Developer Identity heading")
	}
}

func TestRenderGeminiProject(t *testing.T) {
	id := loadTestIdentity(t)
	out, err := Render(id, TargetGeminiProject, &id.Projects[0])
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}
	if !strings.Contains(out, id.Projects[0].Name) {
		t.Error("expected project name in output")
	}
}

func TestRenderGeminiProjectNilError(t *testing.T) {
	id := loadTestIdentity(t)
	_, err := Render(id, TargetGeminiProject, nil)
	if err == nil {
		t.Error("expected error when project is nil")
	}
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
		{TargetCopilot, nil},
		{TargetCline, nil},
		{TargetRooCode, nil},
		{TargetWindsurf, nil},
		{TargetAider, nil},
		{TargetGeminiGlobal, nil},
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

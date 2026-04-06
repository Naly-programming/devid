package scan

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0o644)

	info := DetectProject(dir)
	if !contains(info.Stack, "Go") {
		t.Errorf("expected Go in stack, got %v", info.Stack)
	}
}

func TestDetectTypeScript(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "tsconfig.json"), []byte("{}"), 0o644)

	info := DetectProject(dir)
	if !contains(info.Stack, "TypeScript") {
		t.Errorf("expected TypeScript in stack, got %v", info.Stack)
	}
}

func TestDetectNextJS(t *testing.T) {
	dir := t.TempDir()
	pkg := `{"dependencies":{"next":"14.0.0","react":"18.0.0","typescript":"5.0.0"}}`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkg), 0o644)

	info := DetectProject(dir)
	if !contains(info.Stack, "Next.js") {
		t.Errorf("expected Next.js in stack, got %v", info.Stack)
	}
	if !contains(info.Stack, "TypeScript") {
		t.Errorf("expected TypeScript in stack, got %v", info.Stack)
	}
}

func TestDetectDocker(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM golang"), 0o644)

	info := DetectProject(dir)
	if !contains(info.Infra, "Docker") {
		t.Errorf("expected Docker in infra, got %v", info.Infra)
	}
}

func TestDetectGitHubActions(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".github", "workflows"), 0o755)

	info := DetectProject(dir)
	if !contains(info.Infra, "GitHub Actions") {
		t.Errorf("expected GitHub Actions in infra, got %v", info.Infra)
	}
}

func TestDetectEmpty(t *testing.T) {
	dir := t.TempDir()
	info := DetectProject(dir)
	if len(info.Stack) != 0 {
		t.Errorf("expected empty stack for empty dir, got %v", info.Stack)
	}
}

func TestDetectMultiple(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0o644)
	os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM golang"), 0o644)
	os.MkdirAll(filepath.Join(dir, ".github", "workflows"), 0o755)

	info := DetectProject(dir)
	if !contains(info.Stack, "Go") {
		t.Error("expected Go")
	}
	if !contains(info.Infra, "Docker") {
		t.Error("expected Docker")
	}
	if !contains(info.Infra, "GitHub Actions") {
		t.Error("expected GitHub Actions")
	}
}

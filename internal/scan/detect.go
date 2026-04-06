package scan

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ProjectInfo holds inferred project metadata.
type ProjectInfo struct {
	Stack   []string
	Infra   []string
	Context string
}

// DetectProject scans a directory and infers project stack and context.
func DetectProject(dir string) ProjectInfo {
	var info ProjectInfo

	// Go
	if fileExists(filepath.Join(dir, "go.mod")) {
		info.Stack = append(info.Stack, "Go")
	}

	// Node/TypeScript
	if fileExists(filepath.Join(dir, "package.json")) {
		info.Stack = append(info.Stack, detectFromPackageJSON(dir)...)
	}
	if fileExists(filepath.Join(dir, "tsconfig.json")) {
		if !contains(info.Stack, "TypeScript") {
			info.Stack = append(info.Stack, "TypeScript")
		}
	}

	// Python
	if fileExists(filepath.Join(dir, "requirements.txt")) || fileExists(filepath.Join(dir, "pyproject.toml")) || fileExists(filepath.Join(dir, "setup.py")) {
		info.Stack = append(info.Stack, "Python")
	}

	// Rust
	if fileExists(filepath.Join(dir, "Cargo.toml")) {
		info.Stack = append(info.Stack, "Rust")
	}

	// C# / .NET
	if hasGlob(dir, "*.csproj") || hasGlob(dir, "*.sln") {
		info.Stack = append(info.Stack, "C#", ".NET")
	}

	// Java
	if fileExists(filepath.Join(dir, "pom.xml")) || fileExists(filepath.Join(dir, "build.gradle")) {
		info.Stack = append(info.Stack, "Java")
	}

	// Infra detection
	if fileExists(filepath.Join(dir, "Dockerfile")) || fileExists(filepath.Join(dir, "docker-compose.yml")) || fileExists(filepath.Join(dir, "docker-compose.yaml")) {
		info.Infra = append(info.Infra, "Docker")
	}
	if dirExists(filepath.Join(dir, ".github", "workflows")) {
		info.Infra = append(info.Infra, "GitHub Actions")
	}
	if fileExists(filepath.Join(dir, "vercel.json")) || fileExists(filepath.Join(dir, ".vercel")) {
		info.Infra = append(info.Infra, "Vercel")
	}
	if fileExists(filepath.Join(dir, "terraform.tf")) || dirExists(filepath.Join(dir, "terraform")) {
		info.Infra = append(info.Infra, "Terraform")
	}
	if fileExists(filepath.Join(dir, ".env")) || fileExists(filepath.Join(dir, ".env.local")) {
		// Don't add to infra, but note for context
	}

	// Database detection from common config patterns
	if hasGlob(dir, "*.prisma") || dirExists(filepath.Join(dir, "prisma")) {
		info.Stack = append(info.Stack, "Prisma")
	}
	if dirExists(filepath.Join(dir, "supabase")) || dirExists(filepath.Join(dir, ".supabase")) {
		info.Stack = append(info.Stack, "Supabase")
	}

	return info
}

func detectFromPackageJSON(dir string) []string {
	var stack []string
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return stack
	}

	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return stack
	}

	allDeps := make(map[string]bool)
	for k := range pkg.Dependencies {
		allDeps[k] = true
	}
	for k := range pkg.DevDependencies {
		allDeps[k] = true
	}

	if allDeps["typescript"] {
		stack = append(stack, "TypeScript")
	}
	if allDeps["next"] {
		stack = append(stack, "Next.js")
	} else if allDeps["react"] {
		stack = append(stack, "React")
	}
	if allDeps["vue"] {
		stack = append(stack, "Vue")
	}
	if allDeps["svelte"] || allDeps["@sveltejs/kit"] {
		stack = append(stack, "Svelte")
	}
	if allDeps["express"] {
		stack = append(stack, "Express")
	}
	if allDeps["tailwindcss"] {
		stack = append(stack, "Tailwind")
	}
	if allDeps["playwright"] || allDeps["@playwright/test"] {
		stack = append(stack, "Playwright")
	}
	if allDeps["jest"] {
		stack = append(stack, "Jest")
	}
	if allDeps["vitest"] {
		stack = append(stack, "Vitest")
	}

	if len(stack) == 0 {
		stack = append(stack, "JavaScript")
	}

	return stack
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func hasGlob(dir, pattern string) bool {
	matches, _ := filepath.Glob(filepath.Join(dir, pattern))
	return len(matches) > 0
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

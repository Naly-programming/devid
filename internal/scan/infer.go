package scan

import (
	"os"
	"path/filepath"
	"strings"
)

// InferredSource tracks where an inferred file was found.
type InferredSource struct {
	Path    string
	Content string
}

// FindExistingContextFiles scans common locations for existing AI context files.
// Checks ~/.claude/CLAUDE.md, and scans a base directory for repos with CLAUDE.md,
// AGENTS.md, or .cursor/rules files.
func FindExistingContextFiles(scanDirs []string) []InferredSource {
	var sources []InferredSource

	// Global CLAUDE.md
	home, err := os.UserHomeDir()
	if err == nil {
		globalClaude := filepath.Join(home, ".claude", "CLAUDE.md")
		if content, err := os.ReadFile(globalClaude); err == nil {
			// Strip devid markers if present
			cleaned := stripDevidMarkers(string(content))
			if strings.TrimSpace(cleaned) != "" {
				sources = append(sources, InferredSource{Path: globalClaude, Content: cleaned})
			}
		}
	}

	// Scan provided directories for per-repo context files
	for _, dir := range scanDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			repoDir := filepath.Join(dir, entry.Name())
			for _, name := range []string{"CLAUDE.md", "AGENTS.md", ".cursor/rules/devid.mdc", ".cursorrules"} {
				path := filepath.Join(repoDir, name)
				content, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				cleaned := stripDevidMarkers(string(content))
				if strings.TrimSpace(cleaned) != "" {
					sources = append(sources, InferredSource{Path: path, Content: cleaned})
				}
			}
		}
	}

	return sources
}

// BuildInferencePrompt creates a prompt for an AI to extract identity from existing context files.
func BuildInferencePrompt(sources []InferredSource) string {
	var b strings.Builder

	b.WriteString("The following are existing AI context files found on this developer's machine.\n")
	b.WriteString("Extract a unified developer identity from them.\n")
	b.WriteString("Output valid TOML only, matching the devid schema.\n")
	b.WriteString("Values must be fragments, not sentences. Deduplicate across files.\n")
	b.WriteString("If a preference appears in multiple files, include it once.\n")
	b.WriteString("Omit project-specific details - only extract global identity.\n\n")

	for _, src := range sources {
		b.WriteString("--- File: " + src.Path + " ---\n")
		b.WriteString(src.Content)
		b.WriteString("\n\n")
	}

	return b.String()
}

// stripDevidMarkers removes content between devid markers.
// Returns only content outside the markers (user-written content).
func stripDevidMarkers(content string) string {
	startMarker := "<!-- devid:start -->"
	endMarker := "<!-- devid:end -->"

	startIdx := strings.Index(content, startMarker)
	endIdx := strings.Index(content, endMarker)

	if startIdx >= 0 && endIdx >= 0 {
		before := content[:startIdx]
		after := content[endIdx+len(endMarker):]
		return strings.TrimSpace(before + after)
	}

	return content
}

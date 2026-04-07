package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(diffCmd)
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show what's changed vs distributed files and pending queue",
	RunE:  runDiff,
}

func runDiff(cmd *cobra.Command, args []string) error {
	if !config.Exists() {
		fmt.Println("No identity.toml found. Run `devid init` first.")
		return silentErr{config.ErrNoIdentity}
	}

	id, err := config.Load()
	if err != nil {
		return err
	}

	anyDiff := false

	// Compare identity against distributed targets
	home, _ := os.UserHomeDir()
	targets := []struct {
		name   string
		path   string
		target generate.Target
		proj   *config.Project
	}{
		{"claude-global", filepath.Join(home, ".claude", "CLAUDE.md"), generate.TargetClaudeGlobal, nil},
	}

	for _, t := range targets {
		rendered, err := generate.Render(id, t.target, t.proj)
		if err != nil {
			continue
		}
		wrapped := generate.WrapWithMarkers(rendered)

		existing, err := os.ReadFile(t.path)
		if os.IsNotExist(err) {
			fmt.Printf("Target %s: not yet distributed\n", t.name)
			anyDiff = true
			continue
		}
		if err != nil {
			continue
		}

		// Extract current devid content from the file
		existingStr := string(existing)
		currentContent := extractMarkerContent(existingStr)

		if currentContent != wrapped {
			fmt.Printf("Target %s: out of date\n", t.name)
			showTextDiff(currentContent, wrapped)
			anyDiff = true
		}
	}

	// Queue status
	candidates, _ := devsync.ListQueue()
	if len(candidates) > 0 {
		anyDiff = true
		fmt.Printf("\nPending queue: %d candidates\n", len(candidates))
		for i, c := range candidates {
			fmt.Printf("\n--- Candidate %d (%s, %s) ---\n", i+1, c.Source, c.Timestamp.Format("2006-01-02 15:04"))
			if c.Diff != "" {
				for _, line := range strings.Split(c.Diff, "\n") {
					if strings.HasPrefix(line, "+ ") {
						fmt.Printf("  \033[32m%s\033[0m\n", line)
					} else if strings.HasPrefix(line, "- ") {
						fmt.Printf("  \033[31m%s\033[0m\n", line)
					}
				}
			}
		}
	}

	if !anyDiff {
		fmt.Println("Everything is in sync. No pending changes.")
	}

	return nil
}

// extractMarkerContent returns the full marker block from a file, or the whole file if no markers.
func extractMarkerContent(content string) string {
	startIdx := strings.Index(content, generate.MarkerStart)
	endIdx := strings.Index(content, generate.MarkerEnd)
	if startIdx >= 0 && endIdx >= 0 {
		return content[startIdx : endIdx+len(generate.MarkerEnd)+1] // +1 for trailing newline
	}
	return content
}

// showTextDiff prints a simple line-by-line comparison highlighting differences.
func showTextDiff(old, new string) {
	oldLines := strings.Split(old, "\n")
	newLines := strings.Split(new, "\n")

	// Build a set of old lines for quick lookup
	oldSet := make(map[string]bool)
	for _, l := range oldLines {
		oldSet[l] = true
	}
	newSet := make(map[string]bool)
	for _, l := range newLines {
		newSet[l] = true
	}

	for _, l := range oldLines {
		if !newSet[l] && strings.TrimSpace(l) != "" {
			fmt.Printf("  \033[31m- %s\033[0m\n", l)
		}
	}
	for _, l := range newLines {
		if !oldSet[l] && strings.TrimSpace(l) != "" {
			fmt.Printf("  \033[32m+ %s\033[0m\n", l)
		}
	}
}

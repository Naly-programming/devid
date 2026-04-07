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
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show identity status overview",
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	// Identity file
	if !config.Exists() {
		fmt.Println("No identity configured. Run `devid init` to get started.")
		return nil
	}

	id, err := config.Load()
	if err != nil {
		return err
	}

	idPath, _ := config.IdentityPath()
	info, _ := os.Stat(idPath)

	fmt.Println("Identity")
	fmt.Printf("  name:     %s\n", id.Identity.Name)
	fmt.Printf("  tone:     %s\n", id.Identity.Tone)
	fmt.Printf("  file:     %s\n", idPath)
	if info != nil {
		fmt.Printf("  modified: %s\n", info.ModTime().Format("2006-01-02 15:04"))
	}
	fmt.Printf("  version:  %s\n", id.Meta.Version)

	// Stack summary
	var stack []string
	stack = append(stack, id.Stack.Primary...)
	stack = append(stack, id.Stack.Secondary...)
	if len(stack) > 0 {
		fmt.Printf("  stack:    %s\n", strings.Join(stack, ", "))
	}

	// Projects
	fmt.Printf("\nProjects (%d)\n", len(id.Projects))
	for _, p := range id.Projects {
		fmt.Printf("  %s (%s)\n", p.Name, strings.Join(p.Stack, ", "))
	}
	if len(id.Projects) == 0 {
		fmt.Println("  none - run `devid add` from a repo to add one")
	}

	// Distribution targets
	fmt.Println("\nTargets")
	home, _ := os.UserHomeDir()
	targets := []struct {
		name string
		path string
	}{
		{"claude-global", filepath.Join(home, ".claude", "CLAUDE.md")},
	}

	for _, t := range targets {
		if tInfo, err := os.Stat(t.path); err == nil {
			fmt.Printf("  %-16s %s (updated %s)\n", t.name, t.path, tInfo.ModTime().Format("2006-01-02 15:04"))
		} else {
			fmt.Printf("  %-16s not distributed\n", t.name)
		}
	}

	// Token estimates
	estimates := generate.EstimateAll(id)
	for _, e := range estimates {
		status := "ok"
		if e.Over {
			status = "OVER"
		}
		fmt.Printf("  %-16s ~%d/%d tokens %s\n", e.Target, e.Tokens, e.Budget, status)
	}

	// Queue
	candidates, _ := devsync.ListQueue()
	fmt.Printf("\nQueue: %d pending\n", len(candidates))
	for _, c := range candidates {
		fmt.Printf("  %s from %s\n", c.Timestamp.Format("2006-01-02 15:04"), c.Source)
	}

	// Hook status
	fmt.Println("\nHook")
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	if data, err := os.ReadFile(settingsPath); err == nil && strings.Contains(string(data), "devid hook session-end") {
		fmt.Println("  session-end: installed")
	} else {
		fmt.Println("  session-end: not installed (run `devid hook install`)")
	}

	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		fmt.Println("  api key:     set")
	} else {
		fmt.Println("  api key:     not set (hook needs ANTHROPIC_API_KEY)")
	}

	// Learned entries
	if len(id.Learned.Entries) > 0 {
		fmt.Printf("\nLearned (%d)\n", len(id.Learned.Entries))
		for _, e := range id.Learned.Entries {
			fmt.Printf("  %s\n", e)
		}
	}

	return nil
}

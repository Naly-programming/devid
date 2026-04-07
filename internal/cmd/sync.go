package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/Naly-programming/devid/internal/generate"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func init() {
	syncCmd.Flags().Bool("apply", false, "Read proposed TOML from stdin and queue it")
	syncCmd.Flags().Bool("paste", false, "Read proposed TOML from clipboard and queue it")
	syncCmd.Flags().Bool("now", false, "Apply immediately instead of queuing for review")
	rootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Extract preferences from an AI session",
	RunE:  runSync,
}

func runSync(cmd *cobra.Command, args []string) error {
	now, _ := cmd.Flags().GetBool("now")
	paste, _ := cmd.Flags().GetBool("paste")
	if paste {
		input, err := clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to read clipboard: %w", err)
		}
		if len(input) == 0 {
			return fmt.Errorf("clipboard is empty")
		}
		return syncFromInput(input, now)
	}
	apply, _ := cmd.Flags().GetBool("apply")
	if apply {
		return runSyncApply(now)
	}
	return runSyncPrompt()
}

func runSyncPrompt() error {
	var current *config.Identity
	if config.Exists() {
		var err error
		current, err = config.Load()
		if err != nil {
			return err
		}
	}

	prompt := extract.BuildSyncPrompt(current)
	fmt.Println(prompt)

	if err := clipboard.WriteAll(prompt); err == nil {
		fmt.Println("\n--- Copied to clipboard ---")
	} else {
		fmt.Println("\n--- Copy to clipboard failed, use the output above ---")
	}

	return nil
}

func runSyncApply(now bool) error {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}
	return syncFromInput(string(raw), now)
}

func syncFromInput(input string, applyNow bool) error {
	proposed, err := extract.ParseTOMLResponse(input)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	var current *config.Identity
	if config.Exists() {
		current, err = config.Load()
		if err != nil {
			return err
		}
	} else {
		current = &config.Identity{}
	}

	diff, err := devsync.DiffIdentities(current, proposed)
	if err != nil {
		return fmt.Errorf("failed to compute diff: %w", err)
	}

	hasChanges := false
	for _, line := range splitLines(diff) {
		if len(line) > 0 && (line[0] == '+' || line[0] == '-') {
			hasChanges = true
			break
		}
	}
	if !hasChanges {
		fmt.Println("No changes detected.")
		return nil
	}

	merged := extract.MergeIdentities(current, proposed)

	if applyNow {
		if err := config.Save(merged); err != nil {
			return fmt.Errorf("failed to save: %w", err)
		}
		path, _ := config.IdentityPath()
		fmt.Printf("Identity updated: %s\n", path)
		results := distribute.Distribute(merged)
		printDistributeResults(results)
		fmt.Print(generate.FormatEstimates(generate.EstimateAll(merged)))
		return nil
	}

	candidate := devsync.Candidate{
		Timestamp: time.Now(),
		Source:    "sync",
		Proposed:  proposed,
		Diff:      diff,
	}

	if err := devsync.Enqueue(candidate); err != nil {
		return fmt.Errorf("failed to queue candidate: %w", err)
	}

	fmt.Println("Queued 1 candidate for review. Run `devid review` to approve.")
	fmt.Print(generate.FormatEstimates(generate.EstimateAll(merged)))

	return nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

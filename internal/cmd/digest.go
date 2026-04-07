package cmd

import (
	"fmt"
	"time"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/hook"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	digestCmd.Flags().Int("days", 7, "Number of days to include")
	digestCmd.Flags().Bool("analyze", false, "Send signals to API for identity update suggestions")
	rootCmd.AddCommand(digestCmd)
}

var digestCmd = &cobra.Command{
	Use:   "digest",
	Short: "Summarise what your AI tools learned about you recently",
	Long: `Scans recent Claude Code sessions and shows what preference signals
were detected - corrections, stated preferences, repeated patterns.

  devid digest                # last 7 days
  devid digest --days 30      # last 30 days
  devid digest --analyze      # also call API to suggest identity updates`,
	RunE: runDigest,
}

func runDigest(cmd *cobra.Command, args []string) error {
	days, _ := cmd.Flags().GetInt("days")
	analyze, _ := cmd.Flags().GetBool("analyze")

	report, err := hook.BuildDigest(days)
	if err != nil {
		return err
	}

	fmt.Print(hook.FormatDigest(report))

	if !analyze || report.Signals == 0 {
		return nil
	}

	if !api.Available() {
		fmt.Println("\nSet ANTHROPIC_API_KEY to use --analyze for identity update suggestions.")
		return nil
	}

	// Collect all signal messages across sessions
	var allSignalMessages []hook.Message
	for _, session := range report.SessionsList {
		if session.SignalCount == 0 {
			continue
		}
		messages, err := hook.ReadTranscriptMessages(session.Path)
		if err != nil {
			continue
		}
		// Only grab the signal messages and context
		for i, msg := range messages {
			if msg.Role != "user" {
				continue
			}
			for _, sig := range session.Signals {
				if len(msg.Text) > 120 {
					if msg.Text[:120] == sig[:min(120, len(sig))] {
						if i > 0 {
							allSignalMessages = append(allSignalMessages, messages[i-1])
						}
						allSignalMessages = append(allSignalMessages, msg)
						if i < len(messages)-1 {
							allSignalMessages = append(allSignalMessages, messages[i+1])
						}
					}
				} else if msg.Text == sig {
					if i > 0 {
						allSignalMessages = append(allSignalMessages, messages[i-1])
					}
					allSignalMessages = append(allSignalMessages, msg)
					if i < len(messages)-1 {
						allSignalMessages = append(allSignalMessages, messages[i+1])
					}
				}
			}
		}
	}

	if len(allSignalMessages) == 0 {
		return nil
	}

	// Cap at 40 messages
	if len(allSignalMessages) > 40 {
		allSignalMessages = allSignalMessages[len(allSignalMessages)-40:]
	}

	var current *config.Identity
	if config.Exists() {
		current, _ = config.Load()
	}

	fmt.Println("\nAnalyzing signals for identity updates...")
	proposed, _, err := hook.AnalyzeSession(allSignalMessages, current)
	if err != nil {
		return fmt.Errorf("analysis failed: %w", err)
	}

	if proposed == nil {
		fmt.Println("No novel preferences found - your identity is up to date.")
		return nil
	}

	if current == nil {
		current = &config.Identity{}
	}

	diff, _ := devsync.DiffIdentities(current, proposed)

	candidate := devsync.Candidate{
		Timestamp: time.Now(),
		Source:    "digest",
		Proposed:  proposed,
		Diff:      diff,
	}

	if err := devsync.Enqueue(candidate); err != nil {
		return err
	}

	fmt.Println("Queued 1 candidate for review. Run `devid review` to approve.")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

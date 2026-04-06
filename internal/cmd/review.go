package cmd

import (
	"fmt"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/Naly-programming/devid/internal/review"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(reviewCmd)
}

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review and approve queued identity updates",
	RunE:  runReview,
}

func runReview(cmd *cobra.Command, args []string) error {
	candidates, err := devsync.ListQueue()
	if err != nil {
		return fmt.Errorf("failed to list queue: %w", err)
	}

	if len(candidates) == 0 {
		fmt.Println("No candidates in queue.")
		return nil
	}

	result, err := review.Run(candidates)
	if err != nil {
		return fmt.Errorf("review failed: %w", err)
	}

	// Load current identity
	var current *config.Identity
	if config.Exists() {
		current, err = config.Load()
		if err != nil {
			return err
		}
	} else {
		current = &config.Identity{Meta: config.Meta{Version: "1"}}
	}

	approved := 0
	rejected := 0
	skipped := 0

	for i, decision := range result.Decisions {
		switch decision {
		case review.DecisionApprove:
			current = devsync.ApplyCandidate(current, candidates[i])
			devsync.RemoveCandidate(candidates[i].Timestamp.Unix())
			approved++
		case review.DecisionReject:
			devsync.RemoveCandidate(candidates[i].Timestamp.Unix())
			rejected++
		case review.DecisionSkip:
			skipped++
		default:
			skipped++ // pending = skipped (user quit early)
		}
	}

	if approved > 0 {
		if err := config.Save(current); err != nil {
			return fmt.Errorf("failed to save identity: %w", err)
		}

		results := distribute.Distribute(current)
		fmt.Println()
		printDistributeResults(results)
	}

	fmt.Printf("\nApproved %d, rejected %d, skipped %d.\n", approved, rejected, skipped)
	return nil
}

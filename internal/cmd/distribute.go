package cmd

import (
	"errors"
	"fmt"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(distributeCmd)
}

var distributeCmd = &cobra.Command{
	Use:   "distribute",
	Short: "Render and distribute identity to all targets",
	RunE:  runDistribute,
}

func runDistribute(cmd *cobra.Command, args []string) error {
	id, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNoIdentity) {
			fmt.Println("No identity.toml found. Run `devid init` first.")
			return silentErr{err}
		}
		return fmt.Errorf("failed to load identity: %w", err)
	}

	results := distribute.Distribute(id)
	printDistributeResults(results)
	return nil
}

func printDistributeResults(results []distribute.Result) {
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("  WARN  %-16s %v\n", r.Target, r.Err)
		} else {
			fmt.Printf("  %-8s %-16s %s\n", r.Action, r.Target, r.Path)
		}
	}
}

package cmd

import (
	"errors"
	"fmt"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(snippetCmd)
}

var snippetCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Copy compact identity prompt to clipboard",
	RunE:  runSnippet,
}

func runSnippet(cmd *cobra.Command, args []string) error {
	id, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNoIdentity) {
			fmt.Println("No identity.toml found. Run `devid init` first.")
			return silentErr{err}
		}
		return err
	}

	content, err := generate.Render(id, generate.TargetSnippet, nil)
	if err != nil {
		return fmt.Errorf("failed to render snippet: %w", err)
	}

	if err := clipboard.WriteAll(content); err != nil {
		fmt.Println(content)
		fmt.Printf("\n--- Copy to clipboard failed, snippet printed above (%d chars) ---\n", len(content))
		return nil
	}

	fmt.Printf("Snippet copied to clipboard (%d chars)\n", len(content))
	return nil
}

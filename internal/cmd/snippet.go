package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func init() {
	snippetCmd.Flags().Bool("json", false, "Output as OpenAI system message JSON (for ChatGPT API)")
	snippetCmd.Flags().Bool("print", false, "Print to stdout instead of copying to clipboard")
	rootCmd.AddCommand(snippetCmd)
}

var snippetCmd = &cobra.Command{
	Use:   "snippet",
	Short: "Copy compact identity prompt to clipboard",
	Long: `Outputs a compact identity snippet. Copies to clipboard by default.

Use --json for OpenAI API format (ChatGPT, GPT-4, etc):
  devid snippet --json

Use --print to output to stdout instead of clipboard:
  devid snippet --print`,
	RunE: runSnippet,
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

	asJSON, _ := cmd.Flags().GetBool("json")
	if asJSON {
		msg := map[string]string{
			"role":    "system",
			"content": content,
		}
		out, _ := json.MarshalIndent(msg, "", "  ")
		content = string(out)
	}

	printOnly, _ := cmd.Flags().GetBool("print")
	if printOnly {
		fmt.Print(content)
		return nil
	}

	if err := clipboard.WriteAll(content); err != nil {
		fmt.Print(content)
		fmt.Printf("\n--- Copy to clipboard failed, snippet printed above (%d chars) ---\n", len(content))
		return nil
	}

	fmt.Printf("Snippet copied to clipboard (%d chars)\n", len(content))
	return nil
}

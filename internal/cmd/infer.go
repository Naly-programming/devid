package cmd

import (
	"fmt"
	"os"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/Naly-programming/devid/internal/scan"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func init() {
	inferCmd.Flags().StringSlice("dirs", nil, "Directories to scan for repos (default: parent of cwd)")
	rootCmd.AddCommand(inferCmd)
}

var inferCmd = &cobra.Command{
	Use:   "infer",
	Short: "Infer identity from existing CLAUDE.md and context files",
	Long: `Scans your machine for existing AI context files (CLAUDE.md, AGENTS.md,
.cursor/rules) and extracts a unified identity from them.

If ANTHROPIC_API_KEY is set, sends the files to the API for extraction.
Otherwise, copies the inference prompt to your clipboard for manual use.`,
	RunE: runInfer,
}

func runInfer(cmd *cobra.Command, args []string) error {
	dirs, _ := cmd.Flags().GetStringSlice("dirs")

	// Default to parent of cwd (scan sibling repos)
	if len(dirs) == 0 {
		cwd, err := os.Getwd()
		if err == nil {
			dirs = []string{cwd + "/.."}
		}
	}

	fmt.Println("Scanning for existing context files...")
	sources := scan.FindExistingContextFiles(dirs)

	if len(sources) == 0 {
		fmt.Println("No existing context files found.")
		return nil
	}

	fmt.Printf("Found %d context files:\n", len(sources))
	for _, src := range sources {
		fmt.Printf("  %s\n", src.Path)
	}
	fmt.Println()

	prompt := scan.BuildInferencePrompt(sources)

	if api.Available() {
		fmt.Println("Sending to API for extraction...")
		response, err := api.Call(extract.ExtractionPrompt, prompt)
		if err != nil {
			return fmt.Errorf("API call failed: %w", err)
		}

		return saveInferredIdentity(response)
	}

	// No API key - copy prompt for manual use
	fmt.Println("No ANTHROPIC_API_KEY set. Copying inference prompt to clipboard.")
	fmt.Println("Paste into Claude, then run: devid init --paste")
	fmt.Println()

	if err := clipboard.WriteAll(prompt); err != nil {
		fmt.Println(prompt)
		fmt.Println("\n--- Copy to clipboard failed, prompt printed above ---")
	} else {
		fmt.Println("Prompt copied to clipboard.")
	}

	return nil
}

func saveInferredIdentity(response string) error {
	id, err := extract.ParseTOMLResponse(response)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if id.Meta.Version == "" {
		id.Meta.Version = "1"
	}

	// If identity already exists, merge
	if config.Exists() {
		current, err := config.Load()
		if err == nil {
			id = extract.MergeIdentities(current, id)
		}
	}

	if err := config.Save(id); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	path, _ := config.IdentityPath()
	fmt.Printf("Identity saved to %s\n", path)
	return nil
}

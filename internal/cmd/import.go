package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import an identity from a TOML file or stdin",
	Long: `Import a devid identity.toml as a starting point.
If you already have an identity, imported values are merged on top.

  devid import identity.toml        # from file
  cat identity.toml | devid import  # from stdin`,
	Args: cobra.MaximumNArgs(1),
	RunE: runImport,
}

func runImport(cmd *cobra.Command, args []string) error {
	var input []byte
	var err error

	if len(args) > 0 {
		input, err = os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", args[0], err)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
	}

	incoming, err := extract.ParseTOMLResponse(string(input))
	if err != nil {
		return fmt.Errorf("failed to parse: %w", err)
	}

	if config.Exists() {
		var merge bool
		huh.NewConfirm().
			Title("Identity already exists. Merge imported values on top?").
			Value(&merge).
			Run()

		if !merge {
			fmt.Println("Cancelled.")
			return nil
		}

		current, err := config.Load()
		if err != nil {
			return err
		}
		incoming = extract.MergeIdentities(current, incoming)
	}

	if incoming.Meta.Version == "" {
		incoming.Meta.Version = "1"
	}

	if err := config.Save(incoming); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	path, _ := config.IdentityPath()
	fmt.Printf("Identity imported to %s\n", path)
	return nil
}

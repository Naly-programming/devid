package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/spf13/cobra"
)

func init() {
	exportCmd.Flags().String("out", "", "Output directory (default: stdout)")
	rootCmd.AddCommand(exportCmd)
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export your identity TOML",
	Long: `Exports your identity.toml to stdout or a directory.
Private sections are excluded from the export.

  devid export                  # print to stdout (no private data)
  devid export --out ./backup   # write to directory`,
	RunE: runExport,
}

func runExport(cmd *cobra.Command, args []string) error {
	id, err := config.Load()
	if err != nil {
		if errors.Is(err, config.ErrNoIdentity) {
			fmt.Println("No identity.toml found. Run `devid init` first.")
			return silentErr{err}
		}
		return err
	}

	// Strip private data for export
	clean := id.WithoutPrivate()
	tomlStr, err := extract.FormatIdentityTOML(&clean)
	if err != nil {
		return err
	}

	outDir, _ := cmd.Flags().GetString("out")
	if outDir == "" {
		fmt.Print(tomlStr)
		return nil
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(outDir, "identity.toml")
	if err := os.WriteFile(path, []byte(tomlStr), 0o644); err != nil {
		return err
	}

	fmt.Printf("Exported to %s (private data excluded)\n", path)
	return nil
}

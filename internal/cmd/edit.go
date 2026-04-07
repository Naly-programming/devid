package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(editCmd)
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open identity.toml in your editor",
	RunE:  runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
	if !config.Exists() {
		fmt.Println("No identity.toml found. Run `devid init` first.")
		return silentErr{config.ErrNoIdentity}
	}

	p, err := config.IdentityPath()
	if err != nil {
		return err
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Sensible defaults per platform
		if _, err := exec.LookPath("code"); err == nil {
			editor = "code"
		} else if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else if _, err := exec.LookPath("notepad"); err == nil {
			editor = "notepad"
		} else {
			fmt.Printf("No editor found. Set $EDITOR or open manually:\n  %s\n", p)
			return nil
		}
	}

	c := exec.Command(editor, p)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

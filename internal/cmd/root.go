package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "devid",
	Short:         "Developer identity manager for AI tools",
	Long:          "devid maintains a single source-of-truth developer identity file and distributes it as optimised context to every AI coding tool you use.",
	Version:       "0.2.0",
	SilenceUsage:  true,
	SilenceErrors: true,
}

// silentErr wraps an error that has already been printed to the user.
type silentErr struct{ err error }

func (e silentErr) Error() string { return e.err.Error() }

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		if _, ok := err.(silentErr); !ok {
			fmt.Fprintln(os.Stderr, err)
		}
		return err
	}
	return nil
}

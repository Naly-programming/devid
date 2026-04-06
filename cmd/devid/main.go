package main

import (
	"os"

	cmd "github.com/Naly-programming/devid/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package cmd

import (
	"github.com/Naly-programming/devid/internal/mcp"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(mcpCmd)
}

var mcpCmd = &cobra.Command{
	Use:    "mcp",
	Short:  "Start the MCP server (for Claude.ai and other MCP clients)",
	Hidden: false,
	RunE:   runMCP,
}

func runMCP(cmd *cobra.Command, args []string) error {
	return mcp.Serve()
}

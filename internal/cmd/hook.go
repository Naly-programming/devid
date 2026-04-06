package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/hook"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookSessionEndCmd)
	rootCmd.AddCommand(hookCmd)
}

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage Claude Code session hooks",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the session-end hook into Claude Code settings",
	RunE:  runHookInstall,
}

var hookSessionEndCmd = &cobra.Command{
	Use:    "session-end",
	Short:  "Process a completed session (called by Claude Code hook)",
	Hidden: true,
	RunE:   runHookSessionEnd,
}

func runHookInstall(cmd *cobra.Command, args []string) error {
	if !api.Available() {
		fmt.Println("Warning: ANTHROPIC_API_KEY is not set.")
		fmt.Println("The session-end hook needs an API key to analyze sessions.")
		fmt.Println("Set it in your shell profile: export ANTHROPIC_API_KEY=sk-ant-...")
		fmt.Println()
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	settingsPath := filepath.Join(home, ".claude", "settings.json")

	// Read existing settings
	var settings map[string]any
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]any)
		} else {
			return fmt.Errorf("failed to read settings: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, &settings); err != nil {
			return fmt.Errorf("failed to parse settings: %w", err)
		}
	}

	// Build the hook config
	hookEntry := map[string]any{
		"matcher": "",
		"hooks": []map[string]any{
			{
				"type":    "command",
				"command": "devid hook session-end",
				"timeout": 30,
			},
		},
	}

	// Add to settings under hooks.SessionEnd
	hooks, ok := settings["hooks"].(map[string]any)
	if !ok {
		hooks = make(map[string]any)
	}
	hooks["SessionEnd"] = []any{hookEntry}
	settings["hooks"] = hooks

	// Write back
	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, out, 0o644); err != nil {
		return fmt.Errorf("failed to write settings: %w", err)
	}

	fmt.Printf("Hook installed in %s\n", settingsPath)
	fmt.Println("devid will now analyze sessions when they end.")
	return nil
}

func runHookSessionEnd(cmd *cobra.Command, args []string) error {
	// Read hook input from stdin
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	var hookInput hook.HookInput
	if err := json.Unmarshal(input, &hookInput); err != nil {
		return nil // Silently fail
	}

	// Find the transcript - try stdin field first, then locate by session ID
	transcriptPath := hookInput.TranscriptPath
	if transcriptPath == "" && hookInput.SessionID != "" {
		transcriptPath = hook.FindTranscript(hookInput.SessionID, hookInput.Cwd)
	}
	if transcriptPath == "" {
		return nil // No transcript found, nothing to do
	}

	// Read transcript messages
	messages, err := hook.ReadTranscriptMessages(transcriptPath)
	if err != nil {
		return nil // Silently fail - don't break the user's session close
	}

	if len(messages) == 0 {
		return nil
	}

	// Check for signals before calling the API
	signalCount := hook.CountSignals(messages)
	if signalCount == 0 {
		return nil // No signals, no API call, zero tokens
	}

	// Load current identity
	var current *config.Identity
	if config.Exists() {
		current, err = config.Load()
		if err != nil {
			return nil
		}
	}

	// Analyze the session
	proposed, _, err := hook.AnalyzeSession(messages, current)
	if err != nil {
		return nil // Silently fail
	}
	if proposed == nil {
		return nil // No changes
	}

	// Compute diff
	if current == nil {
		current = &config.Identity{}
	}
	diff, _ := devsync.DiffIdentities(current, proposed)

	// Queue the candidate
	candidate := devsync.Candidate{
		Timestamp: time.Now(),
		Source:    "session-end",
		Proposed:  proposed,
		Diff:      diff,
	}

	return devsync.Enqueue(candidate)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/hook"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	watchCmd.Flags().Bool("once", false, "Scan once and exit (for cron)")
	watchCmd.Flags().Int("interval", 600, "Scan interval in seconds (default 10 minutes)")
	rootCmd.AddCommand(watchCmd)
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Scan recent sessions for identity signals",
	Long: `Periodically scans Claude Code session transcripts for preference
signals. Use --once for a single scan (suitable for cron), or run
without flags for continuous monitoring.

Requires ANTHROPIC_API_KEY to analyze sessions.`,
	RunE: runWatch,
}

func runWatch(cmd *cobra.Command, args []string) error {
	if !api.Available() {
		fmt.Println("ANTHROPIC_API_KEY not set. Watch needs an API key to analyze sessions.")
		return silentErr{fmt.Errorf("no API key")}
	}

	once, _ := cmd.Flags().GetBool("once")
	interval, _ := cmd.Flags().GetInt("interval")

	if once {
		return scanSessions()
	}

	fmt.Printf("Watching for new sessions every %d seconds. Ctrl+C to stop.\n", interval)
	for {
		if err := scanSessions(); err != nil {
			fmt.Printf("Scan error: %v\n", err)
		}
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

// lastScanPath returns the path to the file tracking the last scan timestamp.
func lastScanPath() (string, error) {
	dir, err := config.DevidDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".last_scan"), nil
}

func getLastScanTime() time.Time {
	p, err := lastScanPath()
	if err != nil {
		return time.Time{}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return time.Time{}
	}
	t, err := time.Parse(time.RFC3339, strings.TrimSpace(string(data)))
	if err != nil {
		return time.Time{}
	}
	return t
}

func setLastScanTime(t time.Time) {
	p, err := lastScanPath()
	if err != nil {
		return
	}
	os.WriteFile(p, []byte(t.Format(time.RFC3339)), 0o644)
}

func scanSessions() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	lastScan := getLastScanTime()
	projectsDir := filepath.Join(home, ".claude", "projects")

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil // No projects dir, nothing to scan
	}

	// Find all session JSONL files modified since last scan
	var newSessions []string
	for _, projEntry := range entries {
		if !projEntry.IsDir() {
			continue
		}
		projDir := filepath.Join(projectsDir, projEntry.Name())
		files, err := os.ReadDir(projDir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			info, err := f.Info()
			if err != nil {
				continue
			}
			if info.ModTime().After(lastScan) {
				newSessions = append(newSessions, filepath.Join(projDir, f.Name()))
			}
		}
	}

	if len(newSessions) == 0 {
		fmt.Println("No new sessions since last scan.")
		setLastScanTime(time.Now())
		return nil
	}

	// Sort by modification time, newest first
	sort.Slice(newSessions, func(i, j int) bool {
		iInfo, _ := os.Stat(newSessions[i])
		jInfo, _ := os.Stat(newSessions[j])
		if iInfo == nil || jInfo == nil {
			return false
		}
		return iInfo.ModTime().After(jInfo.ModTime())
	})

	fmt.Printf("Found %d new session(s) to scan.\n", len(newSessions))

	var current *config.Identity
	if config.Exists() {
		current, err = config.Load()
		if err != nil {
			return err
		}
	}

	queued := 0
	for _, sessionPath := range newSessions {
		messages, err := hook.ReadTranscriptMessages(sessionPath)
		if err != nil {
			continue
		}

		signalCount := hook.CountSignals(messages)
		sessionName := filepath.Base(sessionPath)
		fmt.Printf("  %s: %d messages, %d signals", sessionName, len(messages), signalCount)

		if signalCount == 0 {
			fmt.Println(" - skipped")
			continue
		}

		proposed, _, err := hook.AnalyzeSession(messages, current)
		if err != nil {
			fmt.Printf(" - error: %v\n", err)
			continue
		}
		if proposed == nil {
			fmt.Println(" - no changes")
			continue
		}

		if current == nil {
			current = &config.Identity{}
		}
		diff, _ := devsync.DiffIdentities(current, proposed)

		candidate := devsync.Candidate{
			Timestamp: time.Now(),
			Source:    "watch",
			Proposed:  proposed,
			Diff:      diff,
		}

		if err := devsync.Enqueue(candidate); err != nil {
			fmt.Printf(" - queue error: %v\n", err)
			continue
		}

		queued++
		fmt.Println(" - queued")
	}

	setLastScanTime(time.Now())

	if queued > 0 {
		fmt.Printf("\nQueued %d candidate(s). Run `devid review` to approve.\n", queued)
	}

	return nil
}

package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// SessionSummary holds stats for a single session.
type SessionSummary struct {
	Path         string
	Project      string
	ModTime      time.Time
	MessageCount int
	SignalCount  int
	Signals      []string // the actual user messages that matched
}

// DigestReport holds the aggregated weekly digest.
type DigestReport struct {
	Period       string
	Sessions     int
	Messages     int
	Signals      int
	TopSignals   []string
	ByProject    map[string]int
	SessionsList []SessionSummary
}

// BuildDigest scans all sessions from the last N days and builds a report.
func BuildDigest(days int) (*DigestReport, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	projectsDir := filepath.Join(home, ".claude", "projects")

	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return nil, fmt.Errorf("no Claude Code sessions found at %s", projectsDir)
	}

	report := &DigestReport{
		Period:    fmt.Sprintf("last %d days", days),
		ByProject: make(map[string]int),
	}

	for _, projEntry := range entries {
		if !projEntry.IsDir() {
			continue
		}

		projName := projEntry.Name()
		projDir := filepath.Join(projectsDir, projName)
		files, err := os.ReadDir(projDir)
		if err != nil {
			continue
		}

		for _, f := range files {
			if !strings.HasSuffix(f.Name(), ".jsonl") {
				continue
			}
			info, err := f.Info()
			if err != nil || info.ModTime().Before(cutoff) {
				continue
			}

			sessionPath := filepath.Join(projDir, f.Name())
			messages, err := ReadTranscriptMessages(sessionPath)
			if err != nil || len(messages) == 0 {
				continue
			}

			summary := SessionSummary{
				Path:         sessionPath,
				Project:      projName,
				ModTime:      info.ModTime(),
				MessageCount: len(messages),
			}

			// Find signal messages
			for _, msg := range messages {
				if msg.Role != "user" {
					continue
				}
				lower := strings.ToLower(msg.Text)
				for _, kw := range signalKeywords {
					if strings.Contains(lower, kw) {
						summary.SignalCount++
						// Truncate long messages
						text := msg.Text
						if len(text) > 120 {
							text = text[:120] + "..."
						}
						summary.Signals = append(summary.Signals, text)
						break
					}
				}
			}

			report.Sessions++
			report.Messages += summary.MessageCount
			report.Signals += summary.SignalCount
			report.ByProject[projName] += summary.SignalCount
			report.SessionsList = append(report.SessionsList, summary)
		}
	}

	// Sort sessions by mod time, newest first
	sort.Slice(report.SessionsList, func(i, j int) bool {
		return report.SessionsList[i].ModTime.After(report.SessionsList[j].ModTime)
	})

	// Collect top signals across all sessions
	seen := make(map[string]bool)
	for _, s := range report.SessionsList {
		for _, sig := range s.Signals {
			if !seen[sig] && len(report.TopSignals) < 20 {
				report.TopSignals = append(report.TopSignals, sig)
				seen[sig] = true
			}
		}
	}

	return report, nil
}

// FormatDigest produces a human-readable digest report.
func FormatDigest(r *DigestReport) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("devid digest (%s)\n", r.Period))
	b.WriteString(strings.Repeat("-", 50) + "\n\n")

	b.WriteString(fmt.Sprintf("Sessions:  %d\n", r.Sessions))
	b.WriteString(fmt.Sprintf("Messages:  %d\n", r.Messages))
	b.WriteString(fmt.Sprintf("Signals:   %d\n", r.Signals))

	if r.Signals == 0 && r.Sessions > 0 {
		b.WriteString("\nNo preference signals detected. Your identity is holding steady.\n")
		return b.String()
	}

	if r.Sessions == 0 {
		b.WriteString("\nNo sessions found in this period.\n")
		return b.String()
	}

	// By project
	if len(r.ByProject) > 0 {
		b.WriteString("\nSignals by project:\n")
		for proj, count := range r.ByProject {
			if count > 0 {
				b.WriteString(fmt.Sprintf("  %-40s %d signals\n", proj, count))
			}
		}
	}

	// Top signals
	if len(r.TopSignals) > 0 {
		b.WriteString("\nPreference signals found:\n")
		for i, sig := range r.TopSignals {
			b.WriteString(fmt.Sprintf("  %2d. %s\n", i+1, sig))
		}
	}

	// Session breakdown
	if len(r.SessionsList) > 0 {
		b.WriteString("\nSessions with signals:\n")
		for _, s := range r.SessionsList {
			if s.SignalCount > 0 {
				b.WriteString(fmt.Sprintf("  %s  %d msgs, %d signals  (%s)\n",
					s.ModTime.Format("2006-01-02 15:04"),
					s.MessageCount, s.SignalCount,
					s.Project))
			}
		}
	}

	return b.String()
}

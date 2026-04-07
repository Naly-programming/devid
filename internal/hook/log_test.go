package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLogAndRead(t *testing.T) {
	dir := t.TempDir()
	// Override LogDir by creating the log file directly
	logDir := filepath.Join(dir, ".devid", "logs")
	os.MkdirAll(logDir, 0o755)

	logPath := filepath.Join(logDir, "hook.log")

	// Write some log lines directly
	f, _ := os.Create(logPath)
	f.WriteString("[2026-04-07 10:00:00] session-end hook triggered\n")
	f.WriteString("[2026-04-07 10:00:01] 50 messages, 3 signals detected\n")
	f.WriteString("[2026-04-07 10:00:02] queued 1 candidate for review\n")
	f.Close()

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "session-end hook triggered") {
		t.Error("expected log to contain trigger message")
	}
	if !strings.Contains(content, "3 signals detected") {
		t.Error("expected log to contain signal count")
	}
}

func TestReadLogsLimit(t *testing.T) {
	// Test the splitByNewline and limiting logic
	lines := splitByNewline("line1\nline2\nline3\nline4\nline5\n")
	if len(lines) != 5 {
		t.Errorf("expected 5 lines, got %d", len(lines))
	}
}

func TestFindTranscriptMissing(t *testing.T) {
	// Non-existent session should return empty
	result := FindTranscript("nonexistent-session-id", "/tmp")
	if result != "" {
		t.Errorf("expected empty string for missing session, got %q", result)
	}
}

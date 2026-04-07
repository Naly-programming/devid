package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogDir returns the path to ~/.devid/logs/, creating it if needed.
func LogDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".devid", "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// Log writes a timestamped entry to the hook log file.
func Log(format string, args ...any) {
	dir, err := LogDir()
	if err != nil {
		return
	}

	path := filepath.Join(dir, "hook.log")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()

	ts := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(f, "[%s] %s\n", ts, msg)
}

// ReadLogs returns the contents of the hook log file.
func ReadLogs(lines int) (string, error) {
	dir, err := LogDir()
	if err != nil {
		return "", err
	}

	path := filepath.Join(dir, "hook.log")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	if lines <= 0 {
		return string(data), nil
	}

	// Return last N lines
	content := string(data)
	allLines := splitByNewline(content)
	if len(allLines) <= lines {
		return content, nil
	}
	start := len(allLines) - lines
	result := ""
	for i := start; i < len(allLines); i++ {
		result += allLines[i] + "\n"
	}
	return result, nil
}

func splitByNewline(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			if line != "" {
				lines = append(lines, line)
			}
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

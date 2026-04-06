package hook

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// HookInput is the JSON payload Claude Code passes to hooks via stdin.
type HookInput struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	Cwd            string `json:"cwd"`
	HookEventName  string `json:"hook_event_name"`
}

// FindTranscript locates a session's JSONL transcript file.
// Claude Code stores transcripts in ~/.claude/projects/{project-hash}/{session-id}.jsonl
func FindTranscript(sessionID, cwd string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	projectsDir := filepath.Join(home, ".claude", "projects")
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		return ""
	}

	target := sessionID + ".jsonl"

	// If we have a cwd, try the matching project dir first
	if cwd != "" {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			path := filepath.Join(projectsDir, entry.Name(), target)
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	// Fallback: search all project dirs
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		path := filepath.Join(projectsDir, entry.Name(), target)
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// transcriptEntry represents one line of the session JSONL.
type transcriptEntry struct {
	Type    string         `json:"type"`
	Message *entryMessage  `json:"message,omitempty"`
	UUID    string         `json:"uuid,omitempty"`
}

type entryMessage struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

// ReadTranscriptMessages extracts user and assistant text messages from a session JSONL file.
// Returns pairs of (role, text) for messages that contain text content.
func ReadTranscriptMessages(path string) ([]Message, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var messages []Message
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // 1MB line buffer

	for scanner.Scan() {
		var entry transcriptEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue
		}

		if entry.Message == nil {
			continue
		}

		role := entry.Message.Role
		if role != "user" && role != "assistant" {
			continue
		}

		text := extractText(entry.Message.Content)
		if text == "" {
			continue
		}

		messages = append(messages, Message{Role: role, Text: text})
	}

	return messages, scanner.Err()
}

// Message is a simplified role+text pair from a session transcript.
type Message struct {
	Role string
	Text string
}

// extractText pulls plain text out of a message content field.
// Content can be a string or an array of content blocks.
func extractText(raw json.RawMessage) string {
	// Try as string first
	var s string
	if json.Unmarshal(raw, &s) == nil {
		return s
	}

	// Try as array of content blocks
	var blocks []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	}
	if json.Unmarshal(raw, &blocks) == nil {
		var parts []string
		for _, b := range blocks {
			if b.Type == "text" && b.Text != "" {
				parts = append(parts, b.Text)
			}
		}
		return strings.Join(parts, "\n")
	}

	return ""
}

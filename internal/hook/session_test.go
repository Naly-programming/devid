package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadTranscriptMessages(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")

	lines := `{"type":"permission-mode","permissionMode":"default","sessionId":"abc"}
{"type":"user","message":{"role":"user","content":"fix the bug in auth.go"},"uuid":"1"}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"I'll look at auth.go now."}]},"uuid":"2"}
{"type":"user","message":{"role":"user","content":"no, don't use ORM for this - raw SQL only"},"uuid":"3"}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Got it, I'll use raw SQL."}]},"uuid":"4"}
`
	os.WriteFile(path, []byte(lines), 0o644)

	messages, err := ReadTranscriptMessages(path)
	if err != nil {
		t.Fatalf("ReadTranscriptMessages failed: %v", err)
	}

	if len(messages) != 4 {
		t.Fatalf("expected 4 messages, got %d", len(messages))
	}

	if messages[0].Role != "user" || messages[0].Text != "fix the bug in auth.go" {
		t.Errorf("message 0: got role=%q text=%q", messages[0].Role, messages[0].Text)
	}
	if messages[1].Role != "assistant" || messages[1].Text != "I'll look at auth.go now." {
		t.Errorf("message 1: got role=%q text=%q", messages[1].Role, messages[1].Text)
	}
	if messages[2].Role != "user" || messages[2].Text != "no, don't use ORM for this - raw SQL only" {
		t.Errorf("message 2: got role=%q text=%q", messages[2].Role, messages[2].Text)
	}
}

func TestReadTranscriptStringContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")

	// User messages have string content, assistant has array content
	lines := `{"type":"user","message":{"role":"user","content":"I prefer Go over Rust"},"uuid":"1"}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Noted."}]},"uuid":"2"}
`
	os.WriteFile(path, []byte(lines), 0o644)

	messages, err := ReadTranscriptMessages(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(messages))
	}
	if messages[0].Text != "I prefer Go over Rust" {
		t.Errorf("user text = %q", messages[0].Text)
	}
}

func TestReadTranscriptSkipsNonMessages(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "session.jsonl")

	lines := `{"type":"permission-mode","permissionMode":"default"}
{"type":"file-history-snapshot","snapshot":{}}
{"type":"user","message":{"role":"user","content":"hello"},"uuid":"1"}
{"type":"tool-result","message":{"role":"tool","content":"ok"}}
`
	os.WriteFile(path, []byte(lines), 0o644)

	messages, err := ReadTranscriptMessages(path)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
}

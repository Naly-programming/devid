package hook

import (
	"testing"
)

func TestFilterForSignals(t *testing.T) {
	messages := []Message{
		{Role: "assistant", Text: "I'll add the ORM setup."},
		{Role: "user", Text: "no, don't use ORM for this"},
		{Role: "assistant", Text: "Got it, raw SQL."},
		{Role: "user", Text: "looks good"},
		{Role: "assistant", Text: "Done."},
		{Role: "user", Text: "I prefer explicit error types"},
		{Role: "assistant", Text: "Noted."},
	}

	filtered := filterForSignals(messages)

	// "don't" in message 1 should match, pulling in context (0, 1, 2)
	// "prefer" in message 5 should match, pulling in context (4, 5, 6)
	if len(filtered) == 0 {
		t.Fatal("expected filtered messages, got none")
	}

	// Check that the correction message is included
	found := false
	for _, m := range filtered {
		if m.Text == "no, don't use ORM for this" {
			found = true
			break
		}
	}
	if !found {
		t.Error("correction message should be in filtered results")
	}

	// Check that the preference message is included
	found = false
	for _, m := range filtered {
		if m.Text == "I prefer explicit error types" {
			found = true
			break
		}
	}
	if !found {
		t.Error("preference message should be in filtered results")
	}
}

func TestFilterForSignalsNoMatches(t *testing.T) {
	messages := []Message{
		{Role: "user", Text: "fix the bug in auth.go"},
		{Role: "assistant", Text: "Looking at auth.go now."},
		{Role: "user", Text: "thanks, that works"},
	}

	filtered := filterForSignals(messages)
	if len(filtered) != 0 {
		t.Errorf("expected no filtered messages, got %d", len(filtered))
	}
}

func TestFilterForSignalsEmpty(t *testing.T) {
	filtered := filterForSignals(nil)
	if filtered != nil {
		t.Error("expected nil for empty input")
	}
}

func TestCountSignals(t *testing.T) {
	messages := []Message{
		{Role: "user", Text: "don't use that pattern"},
		{Role: "assistant", Text: "ok"},
		{Role: "user", Text: "looks good"},
		{Role: "user", Text: "I prefer Go"},
	}

	count := CountSignals(messages)
	if count != 2 {
		t.Errorf("expected 2 signals, got %d", count)
	}
}

func TestCountSignalsAssistantIgnored(t *testing.T) {
	messages := []Message{
		{Role: "assistant", Text: "I don't think that's right"},
		{Role: "user", Text: "fix the tests"},
	}

	count := CountSignals(messages)
	if count != 0 {
		t.Errorf("expected 0 signals (assistant messages ignored), got %d", count)
	}
}

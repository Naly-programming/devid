package config

import (
	"testing"
)

func TestCheckSensitiveClean(t *testing.T) {
	id := &Identity{
		Identity: IdentitySection{Name: "Nathan", Tone: "direct"},
		Stack:    Stack{Primary: []string{"Go", "TypeScript"}},
	}

	warnings := CheckSensitive(id)
	if len(warnings) != 0 {
		t.Errorf("expected no warnings for clean identity, got %v", warnings)
	}
}

func TestCheckSensitiveAPIKey(t *testing.T) {
	id := &Identity{
		Identity: IdentitySection{Name: "sk-ant-abc123-secret-key"},
	}

	warnings := CheckSensitive(id)
	if len(warnings) == 0 {
		t.Error("expected warning for API key in name field")
	}
}

func TestCheckSensitiveGitHubToken(t *testing.T) {
	id := &Identity{
		Stack: Stack{Primary: []string{"ghp_1234567890abcdef"}},
	}

	warnings := CheckSensitive(id)
	if len(warnings) == 0 {
		t.Error("expected warning for GitHub token in stack")
	}
}

func TestCheckSensitiveJWT(t *testing.T) {
	id := &Identity{
		Learned: Learned{Entries: []string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test"}},
	}

	warnings := CheckSensitive(id)
	if len(warnings) == 0 {
		t.Error("expected warning for JWT in learned entries")
	}
}

func TestCheckSensitivePrivateIgnored(t *testing.T) {
	// Private section should not trigger warnings
	id := &Identity{
		Private: map[string]any{
			"api_key": "sk-ant-secret123",
		},
	}

	warnings := CheckSensitive(id)
	if len(warnings) != 0 {
		t.Errorf("private section should not trigger warnings, got %v", warnings)
	}
}

func TestFormatWarnings(t *testing.T) {
	warnings := []SensitiveWarning{
		{Section: "identity", Field: "name", Value: "sk-ant-abc123..."},
	}

	output := FormatWarnings(warnings)
	if output == "" {
		t.Error("expected non-empty output")
	}

	// Empty warnings should return empty string
	output = FormatWarnings(nil)
	if output != "" {
		t.Error("expected empty output for no warnings")
	}
}

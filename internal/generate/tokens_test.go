package generate

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
)

func TestEstimateTokens(t *testing.T) {
	// 100 chars should be ~25 tokens
	text := "aaaa bbbb cccc dddd eeee ffff gggg hhhh iiii jjjj kkkk llll mmmm nnnn oooo pppp qqqq rrrr ssss tttt"
	tokens := EstimateTokens(text)
	if tokens < 20 || tokens > 30 {
		t.Errorf("expected ~25 tokens for 100 chars, got %d", tokens)
	}
}

func TestEstimateAllWithinBudget(t *testing.T) {
	var id config.Identity
	_, err := toml.DecodeFile("../../schema/identity.toml.example", &id)
	if err != nil {
		t.Fatalf("failed to decode example: %v", err)
	}

	estimates := EstimateAll(&id)
	if len(estimates) == 0 {
		t.Fatal("expected at least one estimate")
	}

	for _, e := range estimates {
		if e.Target == "global" && e.Over {
			t.Errorf("global context over budget: %d tokens (budget: %d)", e.Tokens, e.Budget)
		}
	}
}

func TestFormatEstimates(t *testing.T) {
	estimates := []TokenEstimate{
		{Target: "global", Tokens: 300, Budget: 420, Over: false},
		{Target: "project:test", Tokens: 600, Budget: 500, Over: true},
	}

	output := FormatEstimates(estimates)
	if output == "" {
		t.Error("expected non-empty output")
	}
}

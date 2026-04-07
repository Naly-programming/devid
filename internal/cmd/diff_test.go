package cmd

import (
	"testing"
)

func TestExtractMarkerContent(t *testing.T) {
	input := "# My Notes\n\n<!-- devid:start -->\ngenerated stuff\n<!-- devid:end -->\n\n## Custom\n"
	result := extractMarkerContent(input)

	if result != "<!-- devid:start -->\ngenerated stuff\n<!-- devid:end -->\n" {
		t.Errorf("unexpected result: %q", result)
	}
}

func TestExtractMarkerContentNoMarkers(t *testing.T) {
	input := "# Just a file\nNo markers here.\n"
	result := extractMarkerContent(input)
	if result != input {
		t.Errorf("should return full content when no markers, got %q", result)
	}
}

func TestSplitComma(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"Go, TypeScript, Next.js", 3},
		{"Go", 1},
		{"", 0},
		{"  ,  ,  ", 0},
		{" Go , TypeScript ", 2},
	}

	for _, tc := range tests {
		got := splitComma(tc.input)
		if len(got) != tc.want {
			t.Errorf("splitComma(%q) = %d items, want %d", tc.input, len(got), tc.want)
		}
	}
}

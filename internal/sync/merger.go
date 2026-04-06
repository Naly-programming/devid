package sync

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffIdentities produces a human-readable unified diff between two identities.
func DiffIdentities(current, proposed *config.Identity) (string, error) {
	var curBuf, propBuf bytes.Buffer
	if err := toml.NewEncoder(&curBuf).Encode(current); err != nil {
		return "", err
	}
	if err := toml.NewEncoder(&propBuf).Encode(proposed); err != nil {
		return "", err
	}

	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(curBuf.String(), propBuf.String())
	diffs := dmp.DiffMain(a, b, false)
	diffs = dmp.DiffCharsToLines(diffs, c)
	diffs = dmp.DiffCleanupSemantic(diffs)

	// Build a simple unified-style diff
	var result bytes.Buffer
	for _, d := range diffs {
		switch d.Type {
		case diffmatchpatch.DiffDelete:
			for _, line := range splitLines(d.Text) {
				result.WriteString("- " + line + "\n")
			}
		case diffmatchpatch.DiffInsert:
			for _, line := range splitLines(d.Text) {
				result.WriteString("+ " + line + "\n")
			}
		case diffmatchpatch.DiffEqual:
			for _, line := range splitLines(d.Text) {
				result.WriteString("  " + line + "\n")
			}
		}
	}

	return result.String(), nil
}

// ApplyCandidate merges a candidate's proposed changes onto the current identity.
func ApplyCandidate(current *config.Identity, c Candidate) *config.Identity {
	return extract.MergeIdentities(current, c.Proposed)
}

func splitLines(s string) []string {
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

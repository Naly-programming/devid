package distribute

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
)

// Result describes what happened for a single distribution target.
type Result struct {
	Target string
	Path   string
	Action string // "created", "updated", "unchanged"
	Err    error
}

// repoDetectorFunc can be overridden in tests.
var repoDetectorFunc = detectRepo

// SetRepoDetector overrides repo detection for testing.
func SetRepoDetector(fn func() (root string, name string, err error)) {
	repoDetectorFunc = fn
}

// Distribute renders the identity and writes all distribution targets.
func Distribute(id *config.Identity) []Result {
	var results []Result

	// Global CLAUDE.md
	results = append(results, distributeGlobal(id))

	// Project-scoped targets (only if inside a git repo)
	root, repoName, err := repoDetectorFunc()
	if err == nil {
		proj := matchProject(id, repoName)
		results = append(results, distributeProjectTargets(id, root, proj)...)
	}

	return results
}

func distributeGlobal(id *config.Identity) Result {
	devidDir, err := config.DevidDir()
	if err != nil {
		return Result{Target: "claude-global", Err: fmt.Errorf("cannot find home dir: %w", err)}
	}
	home := filepath.Dir(devidDir) // ~/.devid -> ~

	claudeDir := filepath.Join(home, ".claude")
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return Result{Target: "claude-global", Err: err}
	}

	content, err := generate.Render(id, generate.TargetClaudeGlobal, nil)
	if err != nil {
		return Result{Target: "claude-global", Err: err}
	}

	path := filepath.Join(claudeDir, "CLAUDE.md")
	action, err := writeWithMarkers(path, content)
	return Result{Target: "claude-global", Path: path, Action: action, Err: err}
}

func distributeProjectTargets(id *config.Identity, repoRoot string, proj *config.Project) []Result {
	var results []Result

	// Project CLAUDE.md
	if proj != nil {
		content, err := generate.Render(id, generate.TargetClaudeProject, proj)
		if err != nil {
			results = append(results, Result{Target: "claude-project", Err: err})
		} else {
			path := filepath.Join(repoRoot, "CLAUDE.md")
			action, err := writeWithMarkers(path, content)
			results = append(results, Result{Target: "claude-project", Path: path, Action: action, Err: err})
		}
	}

	// AGENTS.md
	content, err := generate.Render(id, generate.TargetAgentsMD, proj)
	if err != nil {
		results = append(results, Result{Target: "agents-md", Err: err})
	} else {
		path := filepath.Join(repoRoot, "AGENTS.md")
		action, err := writeWithMarkers(path, content)
		results = append(results, Result{Target: "agents-md", Path: path, Action: action, Err: err})
	}

	// Cursor rules (.cursor/rules/devid.mdc)
	content, err = generate.Render(id, generate.TargetCursor, nil)
	if err != nil {
		results = append(results, Result{Target: "cursor", Err: err})
	} else {
		rulesDir := filepath.Join(repoRoot, ".cursor", "rules")
		if err := os.MkdirAll(rulesDir, 0o755); err != nil {
			results = append(results, Result{Target: "cursor", Err: err})
		} else {
			// Cursor .mdc files are standalone - no markers needed, devid owns the whole file
			path := filepath.Join(rulesDir, "devid.mdc")
			action, err := writeFile(path, content)
			results = append(results, Result{Target: "cursor", Path: path, Action: action, Err: err})
		}
	}

	return results
}

// writeWithMarkers writes content between devid markers in a file.
// If the file doesn't exist, creates it with markers.
// If the file exists with markers, replaces content between them.
// If the file exists without markers, prepends the marker block.
// Returns the action taken: "created", "updated", or "unchanged".
func writeWithMarkers(path, content string) (string, error) {
	wrapped := generate.WrapWithMarkers(content)

	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "created", os.WriteFile(path, []byte(wrapped), 0o644)
	}
	if err != nil {
		return "", err
	}

	existingStr := string(existing)
	startIdx := strings.Index(existingStr, generate.MarkerStart)
	endIdx := strings.Index(existingStr, generate.MarkerEnd)

	var newContent string
	if startIdx >= 0 && endIdx >= 0 {
		// Replace content between markers
		after := existingStr[endIdx+len(generate.MarkerEnd):]
		newContent = existingStr[:startIdx] + wrapped + strings.TrimPrefix(after, "\n")
	} else {
		// No markers found - prepend
		newContent = wrapped + "\n" + existingStr
	}

	if newContent == existingStr {
		return "unchanged", nil
	}

	return "updated", os.WriteFile(path, []byte(newContent), 0o644)
}

// writeFile writes content to a file, returning the action taken.
// Unlike writeWithMarkers, this owns the entire file.
func writeFile(path, content string) (string, error) {
	existing, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return "created", os.WriteFile(path, []byte(content), 0o644)
	}
	if err != nil {
		return "", err
	}
	if string(existing) == content {
		return "unchanged", nil
	}
	return "updated", os.WriteFile(path, []byte(content), 0o644)
}

// detectRepo finds the git repo root and extracts the repo name.
func detectRepo() (string, string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", "", fmt.Errorf("not in a git repo: %w", err)
	}
	root := strings.TrimSpace(string(out))
	name := filepath.Base(root)
	return root, name, nil
}

// MatchProject finds the project entry matching the given repo name.
func matchProject(id *config.Identity, repoName string) *config.Project {
	for i := range id.Projects {
		if strings.EqualFold(id.Projects[i].Repo, repoName) {
			return &id.Projects[i]
		}
	}
	return nil
}

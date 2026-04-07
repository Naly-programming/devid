package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/generate"
	devsync "github.com/Naly-programming/devid/internal/sync"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doctorCmd)
}

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose common issues with your devid setup",
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	pass := 0
	warn := 0
	fail := 0

	check := func(name string, ok bool, detail string) {
		if ok {
			pass++
			fmt.Printf("  OK    %s\n", name)
		} else if detail != "" && strings.HasPrefix(detail, "WARN") {
			warn++
			fmt.Printf("  WARN  %s - %s\n", name, detail[5:])
		} else {
			fail++
			fmt.Printf("  FAIL  %s - %s\n", name, detail)
		}
	}

	fmt.Println("devid doctor")
	fmt.Println()

	// 1. Identity exists
	idExists := config.Exists()
	check("identity.toml exists", idExists, "run `devid init` to create one")

	var id *config.Identity
	if idExists {
		var err error
		id, err = config.Load()
		check("identity.toml is valid TOML", err == nil, fmt.Sprintf("%v", err))

		if id != nil {
			// 2. Has a name
			check("identity.name is set", id.Identity.Name != "", "identity has no name")

			// 3. Has stack
			check("stack.primary is set", len(id.Stack.Primary) > 0, "WARN no primary stack defined")

			// 4. Token budget
			estimates := generate.EstimateAll(id)
			for _, e := range estimates {
				if e.Target == "global" {
					check(fmt.Sprintf("token budget (~%d/%d)", e.Tokens, e.Budget), !e.Over,
						fmt.Sprintf("WARN global context is %d tokens, over %d budget", e.Tokens, e.Budget))
				}
			}

			// 5. Sensitive data
			warnings := config.CheckSensitive(id)
			check("no sensitive data in public sections", len(warnings) == 0,
				fmt.Sprintf("WARN %d potential secrets found - run `devid status` for details", len(warnings)))
		}
	}

	// 6. Global CLAUDE.md exists
	home, _ := os.UserHomeDir()
	claudePath := filepath.Join(home, ".claude", "CLAUDE.md")
	_, err := os.Stat(claudePath)
	check("~/.claude/CLAUDE.md exists", err == nil, "run `devid distribute`")

	// 7. Global GEMINI.md exists
	geminiPath := filepath.Join(home, ".gemini", "GEMINI.md")
	_, err = os.Stat(geminiPath)
	check("~/.gemini/GEMINI.md exists", err == nil, "WARN run `devid distribute` to create")

	// 8. Git available
	_, err = exec.LookPath("git")
	check("git is installed", err == nil, "git not found in PATH")

	// 9. In a git repo
	gitOut, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	inRepo := err == nil
	check("inside a git repo", inRepo, "WARN not in a git repo - project targets won't be written")

	// 10. Check repo has distributed targets
	if inRepo {
		repoRoot := strings.TrimSpace(string(gitOut))
		targets := []struct {
			name string
			path string
		}{
			{"AGENTS.md", filepath.Join(repoRoot, "AGENTS.md")},
			{"copilot-instructions.md", filepath.Join(repoRoot, ".github", "copilot-instructions.md")},
		}
		for _, t := range targets {
			_, err := os.Stat(t.path)
			check(fmt.Sprintf("repo has %s", t.name), err == nil,
				fmt.Sprintf("WARN run `devid distribute` in this repo"))
		}
	}

	// 11. Queue status
	candidates, _ := devsync.ListQueue()
	if len(candidates) > 0 {
		check("review queue is empty", false,
			fmt.Sprintf("WARN %d candidates pending - run `devid review`", len(candidates)))
	} else {
		check("review queue is empty", true, "")
	}

	// 12. Hook installed
	settingsPath := filepath.Join(home, ".claude", "settings.json")
	settingsData, err := os.ReadFile(settingsPath)
	hookInstalled := err == nil && strings.Contains(string(settingsData), "devid hook session-end")
	check("session-end hook installed", hookInstalled, "WARN run `devid hook install` for automatic sync")

	// 13. API key
	hasKey := os.Getenv("ANTHROPIC_API_KEY") != ""
	check("ANTHROPIC_API_KEY is set", hasKey, "WARN needed for hook, watch, and API-direct init")

	// Summary
	fmt.Printf("\n%d passed, %d warnings, %d failed\n", pass, warn, fail)
	if fail > 0 {
		return fmt.Errorf("%d checks failed", fail)
	}
	return nil
}

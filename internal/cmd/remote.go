package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/spf13/cobra"
)

func init() {
	remoteCmd.AddCommand(remoteSetCmd)
	remoteCmd.AddCommand(remoteShowCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(pullCmd)
}

var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "Manage sync remote for multi-machine identity",
}

var remoteSetCmd = &cobra.Command{
	Use:   "set <git-url>",
	Short: "Set the git remote for identity sync",
	Long: `Set a git remote URL for syncing your identity across machines.
This initialises ~/.devid/ as a git repo and configures the remote.

  devid remote set git@github.com:you/devid-identity.git
  devid remote set https://github.com/you/devid-identity.git

The remote repo can be private. Create it first on GitHub/GitLab/etc.`,
	Args: cobra.ExactArgs(1),
	RunE: runRemoteSet,
}

var remoteShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current sync remote",
	RunE:  runRemoteShow,
}

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push identity to remote",
	Long:  "Commits any local changes to identity.toml and pushes to the configured remote.",
	RunE:  runPush,
}

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull identity from remote",
	Long:  "Fetches the latest identity from the configured remote and merges it.",
	RunE:  runPull,
}

func devidDir() (string, error) {
	dir, err := config.DevidDir()
	if err != nil {
		return "", err
	}
	return dir, nil
}

func gitInDevid(args ...string) (string, error) {
	dir, err := devidDir()
	if err != nil {
		return "", err
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func isGitRepo() bool {
	_, err := gitInDevid("rev-parse", "--git-dir")
	return err == nil
}

func ensureGitRepo() error {
	if isGitRepo() {
		return nil
	}
	dir, err := devidDir()
	if err != nil {
		return err
	}

	// Init
	if _, err := gitInDevid("init"); err != nil {
		return fmt.Errorf("git init failed: %w", err)
	}

	// Create .gitignore to exclude non-essential files
	gitignore := filepath.Join(dir, ".gitignore")
	os.WriteFile(gitignore, []byte("queue/\nlogs/\n.last_scan\n*.tmp\n"), 0o644)

	// Initial commit
	gitInDevid("add", "-A")
	gitInDevid("commit", "-m", "initial devid identity")

	return nil
}

func runRemoteSet(cmd *cobra.Command, args []string) error {
	url := args[0]

	if err := ensureGitRepo(); err != nil {
		return err
	}

	// Remove existing remote if any
	gitInDevid("remote", "remove", "origin")

	// Add new remote
	if _, err := gitInDevid("remote", "add", "origin", url); err != nil {
		return fmt.Errorf("failed to add remote: %w", err)
	}

	fmt.Printf("Remote set to %s\n", url)
	fmt.Println("Run `devid push` to upload your identity, or `devid pull` to fetch from remote.")
	return nil
}

func runRemoteShow(cmd *cobra.Command, args []string) error {
	if !isGitRepo() {
		fmt.Println("No remote configured. Run `devid remote set <git-url>` first.")
		return nil
	}

	out, err := gitInDevid("remote", "get-url", "origin")
	if err != nil {
		fmt.Println("No remote configured. Run `devid remote set <git-url>` first.")
		return nil
	}

	fmt.Println(out)
	return nil
}

func runPush(cmd *cobra.Command, args []string) error {
	if !isGitRepo() {
		fmt.Println("No remote configured. Run `devid remote set <git-url>` first.")
		return silentErr{fmt.Errorf("no remote")}
	}

	// Stage and commit any changes
	gitInDevid("add", "-A")
	msg := fmt.Sprintf("devid sync %s", time.Now().Format("2006-01-02 15:04"))
	commitOut, err := gitInDevid("commit", "-m", msg)
	if err != nil {
		if strings.Contains(commitOut, "nothing to commit") {
			fmt.Println("Nothing to push - identity is unchanged.")
		}
	}

	// Push
	out, err := gitInDevid("push", "-u", "origin", "main")
	if err != nil {
		// Try master if main doesn't exist
		out, err = gitInDevid("push", "-u", "origin", "master")
		if err != nil {
			return fmt.Errorf("push failed: %s", out)
		}
	}

	fmt.Println("Identity pushed to remote.")
	return nil
}

func runPull(cmd *cobra.Command, args []string) error {
	if !isGitRepo() {
		fmt.Println("No remote configured. Run `devid remote set <git-url>` first.")
		return silentErr{fmt.Errorf("no remote")}
	}

	// Stash any local changes
	gitInDevid("stash")

	// Pull
	out, err := gitInDevid("pull", "origin", "main", "--rebase")
	if err != nil {
		out, err = gitInDevid("pull", "origin", "master", "--rebase")
		if err != nil {
			// Try unstash
			gitInDevid("stash", "pop")
			return fmt.Errorf("pull failed: %s", out)
		}
	}

	// Unstash
	gitInDevid("stash", "pop")

	fmt.Println("Identity pulled from remote.")

	// Redistribute with the updated identity
	if config.Exists() {
		id, err := config.Load()
		if err == nil {
			fmt.Println("Redistributing...")
			results := distribute.Distribute(id)
			printDistributeResults(results)
		}
	}

	return nil
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/Naly-programming/devid/internal/scan"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add [path]",
	Short: "Add a project overlay for a repo",
	Long:  "Scans a repo directory to infer project context and adds a [[projects]] entry to your identity. Defaults to the current directory.",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runAdd,
}

func runAdd(cmd *cobra.Command, args []string) error {
	id, err := config.Load()
	if err != nil {
		fmt.Println("No identity.toml found. Run `devid init` first.")
		return silentErr{err}
	}

	// Resolve repo path
	repoPath := "."
	if len(args) > 0 {
		repoPath = args[0]
	}
	repoPath, err = filepath.Abs(repoPath)
	if err != nil {
		return err
	}

	repoName := filepath.Base(repoPath)

	// Check if project already exists
	for _, p := range id.Projects {
		if strings.EqualFold(p.Repo, repoName) {
			fmt.Printf("Project %q already exists in identity.toml.\n", repoName)
			var overwrite bool
			huh.NewConfirm().Title("Overwrite?").Value(&overwrite).Run()
			if !overwrite {
				return nil
			}
			break
		}
	}

	// Scan the repo
	fmt.Printf("Scanning %s...\n", repoPath)
	detected := scan.DetectProject(repoPath)

	// Pre-fill and let user edit
	var (
		name    = repoName
		repo    = repoName
		stack   = strings.Join(detected.Stack, ", ")
		infra   = strings.Join(detected.Infra, ", ")
		context = detected.Context
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Project name").Value(&name),
			huh.NewInput().Title("Repo name (for matching)").Value(&repo),
			huh.NewInput().Title("Stack (comma-separated)").Value(&stack),
			huh.NewInput().Title("Infra (comma-separated)").Value(&infra),
			huh.NewInput().Title("Context (what is this project?)").Value(&context),
		).Title(fmt.Sprintf("Project: %s", repoName)),
	)

	if err := form.Run(); err != nil {
		return err
	}

	proj := config.Project{
		Name:    name,
		Repo:    repo,
		Stack:   splitComma(stack),
		Infra:   splitComma(infra),
		Context: context,
	}

	// Update or append
	updated := false
	for i, p := range id.Projects {
		if strings.EqualFold(p.Repo, repo) {
			id.Projects[i] = proj
			updated = true
			break
		}
	}
	if !updated {
		id.Projects = append(id.Projects, proj)
	}

	if err := config.Save(id); err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Project %q added to identity.toml\n", name)

	// Distribute if we're in the repo
	cwd, _ := os.Getwd()
	if filepath.Clean(cwd) == filepath.Clean(repoPath) {
		results := distribute.Distribute(id)
		printDistributeResults(results)
	} else {
		fmt.Println("Run `devid distribute` from the repo to generate project-specific files.")
	}

	return nil
}

package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Naly-programming/devid/internal/api"
	"github.com/Naly-programming/devid/internal/config"
	"github.com/Naly-programming/devid/internal/distribute"
	"github.com/Naly-programming/devid/internal/extract"
	"github.com/Naly-programming/devid/internal/generate"
	"github.com/Naly-programming/devid/internal/scan"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func init() {
	initCmd.Flags().Bool("apply", false, "Read identity TOML from stdin (pipe)")
	initCmd.Flags().Bool("paste", false, "Read identity TOML from clipboard")
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Bootstrap your developer identity",
	Long: `Bootstrap your developer identity file.

By default, offers a choice between AI extraction (recommended) and manual form.

  devid init                      # interactive - choose AI or manual
  devid init --paste              # read TOML from clipboard
  devid init --apply < response.toml  # read TOML from stdin`,
	RunE: runInit,
}

func runInit(cmd *cobra.Command, args []string) error {
	paste, _ := cmd.Flags().GetBool("paste")
	if paste {
		return runInitPaste()
	}
	apply, _ := cmd.Flags().GetBool("apply")
	if apply {
		return runInitApply()
	}
	return runInitInteractive()
}

func runInitInteractive() error {
	if config.Exists() {
		var overwrite bool
		err := huh.NewConfirm().
			Title("Identity already exists. Overwrite?").
			Value(&overwrite).
			Run()
		if err != nil {
			return err
		}
		if !overwrite {
			fmt.Println("Cancelled.")
			return nil
		}
	}

	// If API key is available, offer direct API extraction
	options := []huh.Option[string]{
		huh.NewOption("Fill in manually", "manual"),
	}
	if api.Available() {
		options = []huh.Option[string]{
			huh.NewOption("Extract automatically via API (recommended)", "api"),
			huh.NewOption("Extract from AI (copy/paste)", "ai"),
			huh.NewOption("Fill in manually", "manual"),
		}
	} else {
		options = []huh.Option[string]{
			huh.NewOption("Extract from AI (copy/paste)", "ai"),
			huh.NewOption("Fill in manually", "manual"),
		}
	}

	var mode string
	err := huh.NewSelect[string]().
		Title("How do you want to create your identity?").
		Options(options...).
		Value(&mode).
		Run()
	if err != nil {
		return err
	}

	switch mode {
	case "api":
		return runInitAPI()
	case "ai":
		return runInitAI()
	default:
		return runInitManual()
	}
}

func runInitAPI() error {
	fmt.Println("Scanning for existing context files...")

	// Scan for existing CLAUDE.md, .cursorrules, etc
	cwd, _ := os.Getwd()
	sources := scan.FindExistingContextFiles([]string{cwd + "/.."})

	var prompt string
	if len(sources) > 0 {
		fmt.Printf("Found %d existing context files, using them as input.\n", len(sources))
		prompt = scan.BuildInferencePrompt(sources)
	} else {
		fmt.Println("No existing context files found, using blank extraction prompt.")
		prompt = extract.BuildSyncPrompt(nil)
	}

	fmt.Println("Calling Claude API...")
	response, err := api.Call(extract.ExtractionPrompt, prompt)
	if err != nil {
		return fmt.Errorf("API call failed: %w", err)
	}

	return saveFromResponse(response)
}

func runInitAI() error {
	prompt := extract.BuildSyncPrompt(nil)
	fmt.Println(prompt)

	if err := clipboard.WriteAll(prompt); err == nil {
		fmt.Println("\n--- Copied to clipboard ---")
	} else {
		fmt.Println("\n--- Copy to clipboard failed, use the output above ---")
	}

	fmt.Println()
	fmt.Println("Paste this prompt into Claude (claude.ai, Claude Code, or any AI tool).")
	fmt.Println("Copy the TOML response, then run:")
	fmt.Println()
	fmt.Println("  devid init --paste                                  # read from clipboard")
	fmt.Println()

	return nil
}

func runInitPaste() error {
	input, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read clipboard: %w", err)
	}
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("clipboard is empty")
	}
	return saveFromResponse(input)
}

func runInitApply() error {
	raw, err := io.ReadAll(os.Stdin)
	if err != nil {
		return fmt.Errorf("failed to read stdin: %w", err)
	}
	return saveFromResponse(string(raw))
}

func saveFromResponse(input string) error {
	id, err := extract.ParseTOMLResponse(input)
	if err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if id.Meta.Version == "" {
		id.Meta.Version = "1"
	}

	if err := config.Save(id); err != nil {
		return fmt.Errorf("failed to save identity: %w", err)
	}

	path, _ := config.IdentityPath()
	fmt.Printf("Identity saved to %s\n", path)

	results := distribute.Distribute(id)
	printDistributeResults(results)
	fmt.Print(generate.FormatEstimates(generate.EstimateAll(id)))

	return nil
}

func runInitManual() error {
	var (
		name      string
		tone      string
		comments  string
		responses string
		pace      string

		primaryStr   string
		secondaryStr string
		dataStr      string
		infraStr     string
		avoidStr     string

		formatting    string
		prStyle       string
		commitStyle   string
		errorHandling string
		naming        string

		verbosity    string
		confirmation string
		suggestions  string
		codeComments string
		tests        string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Your name").
				Value(&name),
			huh.NewInput().
				Title("Tone (how you communicate, 5 words max)").
				Placeholder("direct, plain-spoken, no fluff").
				Value(&tone),
			huh.NewInput().
				Title("Code comment style").
				Placeholder("sound like the dev wrote it, not a textbook").
				Value(&comments),
			huh.NewInput().
				Title("Response format preference").
				Placeholder("prose over bullets, hyphens not em dashes").
				Value(&responses),
			huh.NewInput().
				Title("Pace preference").
				Placeholder("move fast, skip obvious explanations").
				Value(&pace),
		).Title("Identity"),

		huh.NewGroup(
			huh.NewInput().
				Title("Primary languages/frameworks (comma-separated)").
				Placeholder("Go, TypeScript, Next.js").
				Value(&primaryStr),
			huh.NewInput().
				Title("Secondary languages/frameworks (comma-separated)").
				Placeholder("C#, .NET").
				Value(&secondaryStr),
			huh.NewInput().
				Title("Data tools (comma-separated)").
				Placeholder("PostgreSQL, Supabase").
				Value(&dataStr),
			huh.NewInput().
				Title("Infra tools (comma-separated)").
				Placeholder("Docker, GitHub Actions, Vercel").
				Value(&infraStr),
			huh.NewInput().
				Title("Things to avoid (comma-separated)").
				Description("Tools/patterns you explicitly do not want suggested").
				Placeholder("Prisma, ORM abstraction over raw SQL").
				Value(&avoidStr),
		).Title("Stack"),

		huh.NewGroup(
			huh.NewInput().
				Title("Formatting preferences").
				Placeholder("hyphens not em dashes").
				Value(&formatting),
			huh.NewInput().
				Title("PR style").
				Placeholder("small focused PRs, one concern per PR").
				Value(&prStyle),
			huh.NewInput().
				Title("Commit style").
				Placeholder("conventional commits, lowercase, imperative").
				Value(&commitStyle),
			huh.NewInput().
				Title("Error handling approach").
				Placeholder("explicit, no silent swallows, log with context").
				Value(&errorHandling),
			huh.NewInput().
				Title("Naming conventions").
				Placeholder("clear over clever, full words not abbreviations").
				Value(&naming),
		).Title("Conventions"),

		huh.NewGroup(
			huh.NewInput().
				Title("Verbosity preference").
				Placeholder("concise, skip preamble, get to the point").
				Value(&verbosity),
			huh.NewInput().
				Title("Confirmation style").
				Placeholder("don't ask permission for obvious next steps").
				Value(&confirmation),
			huh.NewInput().
				Title("Suggestion style").
				Placeholder("challenge assumptions, flag alternatives").
				Value(&suggestions),
			huh.NewInput().
				Title("Code comments preference").
				Placeholder("minimal, only non-obvious logic").
				Value(&codeComments),
			huh.NewInput().
				Title("Test writing preference").
				Placeholder("write them, don't ask if I want them").
				Value(&tests),
		).Title("AI Preferences"),
	)

	if err := form.Run(); err != nil {
		return err
	}

	id := &config.Identity{
		Meta: config.Meta{Version: "1"},
		Identity: config.IdentitySection{
			Name:      name,
			Tone:      tone,
			Comments:  comments,
			Responses: responses,
			Pace:      pace,
		},
		Stack: config.Stack{
			Primary:   splitComma(primaryStr),
			Secondary: splitComma(secondaryStr),
			Data:      splitComma(dataStr),
			Infra:     splitComma(infraStr),
			Avoid: config.StackAvoid{
				Items: splitComma(avoidStr),
			},
		},
		Conventions: config.Conventions{
			Formatting:    formatting,
			PRStyle:       prStyle,
			CommitStyle:   commitStyle,
			ErrorHandling: errorHandling,
			Naming:        naming,
		},
		AI: config.AI{
			Verbosity:    verbosity,
			Confirmation: confirmation,
			Suggestions:  suggestions,
			CodeComments: codeComments,
			Tests:        tests,
		},
	}

	if err := config.Save(id); err != nil {
		return fmt.Errorf("failed to save identity: %w", err)
	}

	path, _ := config.IdentityPath()
	fmt.Printf("Identity saved to %s\n", path)

	var dist bool
	err := huh.NewConfirm().
		Title("Distribute now?").
		Affirmative("Yes").
		Negative("No").
		Value(&dist).
		Run()
	if err != nil {
		return err
	}

	if dist {
		results := distribute.Distribute(id)
		printDistributeResults(results)
	}

	return nil
}

func splitComma(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

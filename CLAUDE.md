# devid

Go CLI tool that maintains a developer identity TOML file and distributes it as context to AI coding tools.

## Stack
Go, cobra, huh (charmbracelet), bubbletea, lipgloss, BurntSushi/toml

## Key decisions
- No Viper - BurntSushi/toml handles the single TOML file directly
- Section markers (`<!-- devid:start -->` / `<!-- devid:end -->`) for distribution targets - devid owns content between markers, leaves everything else untouched
- `[private]` section in identity.toml is excluded from all rendered output via `Identity.WithoutPrivate()`
- Atomic writes via temp file + os.Rename to prevent corruption
- Queue-based sync: candidates are individual TOML files in `~/.devid/queue/`

## Conventions
- Conventional commits, lowercase, imperative mood
- Explicit error handling, no silent swallows
- Golden file tests in `testdata/golden/` for renderer output (run tests with `-update` to regenerate)
- Test injection via `config.SetHomeDir()` and `distribute.SetRepoDetector()` to avoid touching real filesystem

## Structure
- `cmd/devid/` - entrypoint
- `internal/cmd/` - cobra commands
- `internal/config/` - TOML schema and load/save
- `internal/extract/` - extraction prompts and TOML response parsing
- `internal/generate/` - renderers for each distribution target
- `internal/distribute/` - file writing with marker management
- `internal/sync/` - candidate queue and diff/merge
- `internal/review/` - bubbletea TUI for approving/rejecting candidates

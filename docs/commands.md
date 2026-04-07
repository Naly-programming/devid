# Commands

## Identity management

| Command | Description |
|---------|-------------|
| `devid init` | Bootstrap your identity - API extraction, copy/paste, or manual form |
| `devid init --paste` | Create identity from TOML on your clipboard |
| `devid init --apply` | Create identity from TOML piped via stdin |
| `devid edit` | Open identity.toml in your editor |
| `devid export` | Export identity to stdout (private data excluded) |
| `devid import [file]` | Import an identity from a TOML file or stdin |

## Distribution

| Command | Description |
|---------|-------------|
| `devid distribute` | Render and write identity to all target files |
| `devid snippet` | Copy compact identity to clipboard |
| `devid snippet --json` | Output as OpenAI system message JSON |
| `devid snippet --print` | Print to stdout instead of clipboard |
| `devid add [path]` | Scan a repo and add a project overlay |
| `devid mcp` | Start the MCP server for Claude.ai and other MCP clients |

## Sync and review

| Command | Description |
|---------|-------------|
| `devid sync` | Print extraction prompt, copy to clipboard |
| `devid sync --paste` | Read AI response from clipboard and queue for review |
| `devid sync --apply` | Pipe AI response from stdin and queue for review |
| `devid sync --now` | Combine with --paste or --apply to skip the review queue |
| `devid review` | Approve/reject queued identity updates (TUI) |
| `devid infer` | Infer identity from existing CLAUDE.md files across repos |

## Multi-machine

| Command | Description |
|---------|-------------|
| `devid remote set <url>` | Set git remote for syncing identity across machines |
| `devid remote show` | Show current sync remote |
| `devid push` | Commit and push identity to remote |
| `devid pull` | Fetch identity from remote and redistribute |

## Monitoring

| Command | Description |
|---------|-------------|
| `devid status` | Show identity overview - targets, tokens, queue, hook status |
| `devid diff` | Show what's changed vs distributed files and pending queue |
| `devid digest` | Summarise what AI tools learned about you recently |
| `devid digest --analyze` | Also suggest identity updates via API |
| `devid doctor` | Run diagnostic checks on your setup |
| `devid update` | Update devid to the latest version |
| `devid hook install` | Wire up automatic session-end analysis in Claude Code |
| `devid hook logs` | Show recent hook activity |
| `devid watch` | Scan recent sessions for identity signals |
| `devid watch --once` | Single scan, then exit (for cron) |

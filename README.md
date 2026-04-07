# devid

Developer identity manager for AI tools. Maintains a single source-of-truth identity file (TOML) and distributes it as optimised context to Claude Code, Cursor, Claude.ai, and anything else that accepts a context file.

The core problem: every AI tool starts each session knowing nothing about you. You repeat yourself constantly - your stack, your tone, your conventions. devid captures that once and keeps it everywhere, automatically.

https://github.com/user-attachments/assets/879e222f-2962-4cf0-8faf-76493c7a9f15

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Naly-programming/devid/main/install.sh | sh
```

Or with Go:

```bash
go install github.com/Naly-programming/devid/cmd/devid@latest
```

## Quick start

```bash
# Bootstrap your identity (interactive - choose API, copy/paste, or manual)
devid init

# Or let AI extract your identity, then paste the response back
devid init --paste

# Distribute to all AI tools
devid distribute

# Check everything is in sync
devid status

# Install automatic session sync (optional, needs ANTHROPIC_API_KEY)
devid hook install
```

## Commands

| Command | Description |
|---------|-------------|
| `devid init` | Bootstrap your identity - API extraction, copy/paste, or manual form |
| `devid init --paste` | Create identity from TOML on your clipboard |
| `devid init --apply` | Create identity from TOML piped via stdin |
| `devid distribute` | Render and write identity to all target files |
| `devid status` | Show identity overview - targets, tokens, queue, hook status |
| `devid diff` | Show what's changed vs distributed files and pending queue |
| `devid edit` | Open identity.toml in your editor |
| `devid sync` | Print extraction prompt for updating an existing identity |
| `devid sync --paste` | Read AI response from clipboard and queue for review |
| `devid sync --apply` | Pipe AI response to queue a candidate update |
| `devid sync --now` | Combine with --paste or --apply to skip the review queue |
| `devid review` | Approve/reject queued identity updates (TUI) |
| `devid snippet` | Copy compact identity to clipboard (for claude.ai) |
| `devid add [path]` | Scan a repo and add a project overlay to your identity |
| `devid infer` | Infer identity from existing CLAUDE.md files across your repos |
| `devid watch` | Scan recent sessions for identity signals (--once for cron) |
| `devid mcp` | Start the MCP server for Claude.ai and other MCP clients |
| `devid hook install` | Wire up automatic session-end analysis in Claude Code |
| `devid hook logs` | Show recent hook activity for debugging |

## How it works

1. Your identity lives in `~/.devid/identity.toml` - a compressed TOML file capturing your tone, stack, conventions, and preferences
2. `devid distribute` renders target-specific versions and writes them to:
   - `~/.claude/CLAUDE.md` (global Claude Code context)
   - `{repo}/CLAUDE.md` (per-project overlay, if a matching project is configured)
   - `{repo}/AGENTS.md` (cross-tool compatibility)
   - `{repo}/.cursor/rules/devid.mdc` (Cursor, with YAML frontmatter)
3. Content is wrapped in `<!-- devid:start -->` / `<!-- devid:end -->` markers so your own notes in these files are preserved
4. Optionally, the session-end hook or watch command monitors Claude Code sessions for corrections and preferences, queuing them for review

## AI extraction flow

The recommended way to create your identity - let an AI that already knows you do the work:

```bash
# With API key (one command, zero copy-paste)
export ANTHROPIC_API_KEY=sk-ant-...
devid init                # picks "Extract automatically via API"
                          # scans existing context files, calls Claude, done

# Without API key (copy-paste flow)
devid init                # picks "Extract from AI" - prompt copied to clipboard
                          # paste into Claude, copy the TOML response
devid init --paste        # reads clipboard, saves identity, distributes

# Updating an existing identity
devid sync                # extraction prompt copied to clipboard
devid sync --paste        # queue for review
devid sync --paste --now  # apply immediately, skip queue

devid review              # approve/reject queued changes in TUI
```

## Automatic session sync

devid can automatically analyze your Claude Code sessions, picking up corrections and preference changes without any manual effort. Two approaches:

**Session-end hook** (fires when a Claude Code session closes):

```bash
export ANTHROPIC_API_KEY=sk-ant-...
devid hook install            # adds SessionEnd hook to ~/.claude/settings.json
devid hook logs               # check what the hook has been doing
```

**Watch mode** (scans all recent sessions):

```bash
devid watch --once            # scan once and exit (for cron)
devid watch --interval 300    # continuous, scan every 5 minutes
```

**How it stays token-efficient:** devid pre-filters session transcripts for high-signal keywords - corrections ("don't", "stop", "no"), preferences ("prefer", "always", "instead"), and style instructions ("be more", "be less"). If no signals are found in a session, no API call is made. Zero tokens spent on sessions where nothing identity-relevant happened.

## MCP server

devid can run as an MCP server, giving Claude.ai (or any MCP client) direct access to your identity without clipboard pasting.

```bash
devid mcp    # starts JSON-RPC server on stdin/stdout
```

Available tools: `get_identity`, `get_snippet`, `get_project`. Configure your MCP client to run `devid mcp` as a stdio transport.

## Visibility

```bash
# Full overview of your identity, targets, tokens, queue, and hook status
devid status

# See what's out of date or pending
devid diff

# Quick edit
devid edit
```

`devid distribute`, `devid init`, and `devid sync` all show token estimates so you can verify your identity stays within budget (~400 tokens for global context). Sensitive data in non-private sections triggers a warning on save.

## Project overlays

Add per-project context so AI tools know the specifics of each repo:

```bash
# From inside a repo - scans go.mod, package.json, Dockerfile, etc
devid add

# Or point to a specific path
devid add ~/projects/myapp
```

devid detects your stack (Go, TypeScript, Next.js, React, Python, Rust, etc) and infra (Docker, GitHub Actions, Vercel) automatically. You review and edit before saving.

## Inferring from existing files

Already have CLAUDE.md or .cursor/rules files scattered across repos? devid can scan them and extract a unified identity:

```bash
devid infer                           # scans sibling repos by default
devid infer --dirs ~/projects,~/work  # scan specific directories
```

If `ANTHROPIC_API_KEY` is set, it sends the files to the API for extraction. Otherwise it copies a prompt to your clipboard for manual use.

## Identity schema

The TOML schema is designed for maximum signal per token. Values are fragments, not sentences.

```toml
[identity]
name = "Nathan"
tone = "direct, plain-spoken, no fluff, northern"
comments = "sound like the dev wrote it, not a textbook"

[stack]
primary = ["Go", "TypeScript", "Next.js"]
data = ["PostgreSQL", "Supabase"]

[stack.avoid]
items = ["Prisma", "ORM abstraction over raw SQL"]

[conventions]
pr_style = "small focused PRs, one concern per PR"
commit_style = "conventional commits, lowercase, imperative mood"

[ai]
verbosity = "concise, skip preamble, get to the point"
tests = "write them, dont ask if I want them"

[private]
# Fields here are never included in distributed output
api_key = "..."
```

See `schema/identity.toml.example` for the full annotated schema.

## Distribution targets

| Target | Path | When |
|--------|------|------|
| Claude Code (global) | `~/.claude/CLAUDE.md` | Always |
| Claude Code (project) | `{repo}/CLAUDE.md` | When repo matches a `[[projects]]` entry |
| AGENTS.md | `{repo}/AGENTS.md` | When in a git repo |
| Cursor | `{repo}/.cursor/rules/devid.mdc` | When in a git repo |
| Clipboard snippet | `devid snippet` | On demand (for claude.ai) |
| MCP server | `devid mcp` | On demand (for any MCP client) |

## File locations

```
~/.devid/
  identity.toml        # source of truth
  queue/               # pending candidate updates
  logs/                # hook activity logs
  .last_scan           # watch timestamp tracker

~/.claude/
  settings.json        # hook config (after devid hook install)
  CLAUDE.md            # global identity (after devid distribute)
```

## License

MIT

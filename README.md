# devid

Developer identity manager for AI tools. Maintains a single source-of-truth identity file (TOML) and distributes it as optimised context to Claude Code, Cursor, Claude.ai, and anything else that accepts a context file.

The core problem: every AI tool starts each session knowing nothing about you. You repeat yourself constantly - your stack, your tone, your conventions. devid captures that once and keeps it everywhere, automatically.

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
# Bootstrap your identity (interactive - choose AI extraction or manual)
devid init

# Or let AI extract your identity, then paste the response back
devid init --paste

# Distribute to all AI tools
devid distribute

# Install automatic session sync (optional, needs ANTHROPIC_API_KEY)
devid hook install
```

## Commands

| Command | Description |
|---------|-------------|
| `devid init` | Bootstrap your identity - extract from AI or fill in manually |
| `devid init --paste` | Create identity from TOML on your clipboard |
| `devid init --apply` | Create identity from TOML piped via stdin |
| `devid distribute` | Render and write identity to all target files |
| `devid sync` | Print extraction prompt for updating an existing identity |
| `devid sync --apply` | Pipe AI response to queue a candidate update |
| `devid review` | Approve/reject queued identity updates (TUI) |
| `devid snippet` | Copy compact identity to clipboard (for claude.ai) |
| `devid add [path]` | Scan a repo and add a project overlay to your identity |
| `devid infer` | Infer identity from existing CLAUDE.md files across your repos |
| `devid hook install` | Wire up automatic session-end analysis in Claude Code |

## How it works

1. Your identity lives in `~/.devid/identity.toml` - a compressed TOML file capturing your tone, stack, conventions, and preferences
2. `devid distribute` renders target-specific versions and writes them to:
   - `~/.claude/CLAUDE.md` (global Claude Code context)
   - `{repo}/CLAUDE.md` (per-project overlay, if a matching project is configured)
   - `{repo}/AGENTS.md` (cross-tool compatibility)
   - `{repo}/.cursor/rules` (Cursor)
3. Content is wrapped in `<!-- devid:start -->` / `<!-- devid:end -->` markers so your own notes in these files are preserved
4. Optionally, the session-end hook monitors Claude Code sessions for corrections and preferences, queuing them for review

## AI extraction flow

The recommended way to create your identity - let an AI that already knows you do the work:

```bash
# First time setup
devid init                # choose "Extract from AI" - prompt copied to clipboard
                          # paste into Claude, copy the TOML response
devid init --paste        # reads clipboard, saves identity, distributes

# Updating an existing identity
devid sync                # extraction prompt copied to clipboard
                          # paste into Claude, copy the TOML response
devid sync --apply        # or pipe: powershell Get-Clipboard | devid sync --apply

devid review              # approve/reject changes in TUI
```

## Automatic session sync

devid can automatically analyze your Claude Code sessions when they end, picking up corrections and preference changes without any manual effort.

```bash
# Set your API key (needed for session analysis)
export ANTHROPIC_API_KEY=sk-ant-...

# Install the session-end hook
devid hook install

# That's it - devid now runs silently at session end
# If it finds preference signals, they're queued for review
devid review
```

**How it stays token-efficient:** devid pre-filters session transcripts for high-signal keywords - corrections ("don't", "stop", "no"), preferences ("prefer", "always", "instead"), and style instructions ("be more", "be less"). If no signals are found in a session, no API call is made. Zero tokens spent on sessions where nothing identity-relevant happened.

When signals are found, only the matching messages and their surrounding context are sent to the API with a focused diff prompt that includes your current identity. The API only returns new or changed fields - not a full re-extraction.

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

`devid distribute` shows token estimates after writing, so you can verify your identity stays within budget (~400 tokens for global context). Sensitive data in non-private sections triggers a warning on save.

## Distribution targets

| Target | Path | When |
|--------|------|------|
| Claude Code (global) | `~/.claude/CLAUDE.md` | Always |
| Claude Code (project) | `{repo}/CLAUDE.md` | When repo matches a `[[projects]]` entry |
| AGENTS.md | `{repo}/AGENTS.md` | When in a git repo |
| Cursor | `{repo}/.cursor/rules` | When in a git repo |
| Clipboard snippet | `devid snippet` | On demand (for claude.ai) |

## File locations

```
~/.devid/
  identity.toml        # source of truth
  queue/               # pending candidate updates

~/.claude/
  settings.json        # hook config (after devid hook install)
  CLAUDE.md            # global identity (after devid distribute)
```

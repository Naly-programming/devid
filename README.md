# devid

Developer identity manager for AI tools. Maintains a single source-of-truth identity file (TOML) and distributes it as optimised context to Claude Code, Cursor, Claude.ai, and anything else that accepts a context file.

The core problem: every AI tool starts each session knowing nothing about you. You repeat yourself constantly - your stack, your tone, your conventions. devid captures that once and keeps it everywhere, automatically.

## Install

```
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

# Copy a compact snippet to clipboard (for claude.ai)
devid snippet
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
| `devid snippet` | Copy compact identity to clipboard |

## How it works

1. Your identity lives in `~/.devid/identity.toml` - a compressed TOML file capturing your tone, stack, conventions, and preferences
2. `devid distribute` renders target-specific versions and writes them to:
   - `~/.claude/CLAUDE.md` (global Claude Code context)
   - `{repo}/CLAUDE.md` (per-project overlay, if a matching project is configured)
   - `{repo}/AGENTS.md` (cross-tool compatibility)
   - `{repo}/.cursor/rules` (Cursor)
3. Content is wrapped in `<!-- devid:start -->` / `<!-- devid:end -->` markers so your own notes in these files are preserved

## AI extraction flow

The recommended way to create or update your identity:

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
```

See `schema/identity.toml.example` for the full annotated schema.

# devid

Developer identity manager for AI tools.

One TOML file. 13 distribution targets. Every AI coding tool knows who you are from the first message.

https://github.com/user-attachments/assets/52675d0c-2410-41b8-8226-3bead9224974

---

## The problem

Every AI tool starts each session knowing nothing about you. You repeat the same things constantly - your stack, your tone, your conventions, your preferences. Across Claude Code, Cursor, Copilot, Gemini, Cline, Windsurf, Aider, and everything else.

devid captures that once and keeps it everywhere, automatically.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/Naly-programming/devid/main/install.sh | sh
```

Also available via `npm install -g devid-cli` or `go install github.com/Naly-programming/devid/cmd/devid@latest`.

## Quick start

```bash
devid init              # create your identity (AI extraction or manual)
devid distribute        # push to all 13 targets
devid status            # check everything is in sync
```

## What it distributes to

| Tool | File |
|------|------|
| Claude Code | `~/.claude/CLAUDE.md` and `{repo}/CLAUDE.md` |
| Gemini | `~/.gemini/GEMINI.md` and `{repo}/GEMINI.md` |
| GitHub Copilot | `{repo}/.github/copilot-instructions.md` |
| Cursor | `{repo}/.cursor/rules/devid.mdc` |
| Cline | `{repo}/.clinerules` |
| Roo Code | `{repo}/.roo/rules/devid.md` |
| Windsurf | `{repo}/.windsurf/rules/devid.md` |
| Aider | `{repo}/CONVENTIONS.md` |
| AGENTS.md | `{repo}/AGENTS.md` |
| ChatGPT | `devid snippet` (clipboard) or `devid snippet --json` (API) |
| Any MCP client | `devid mcp` (JSON-RPC server) |

One command. All targets. Section markers preserve any notes you've already got in those files.

## The identity file

`~/.devid/identity.toml` - fragments, not sentences. Maximum signal per token.

```toml
[identity]
name = "Nathan"
tone = "direct, plain-spoken, no fluff, northern"

[stack]
primary = ["Go", "TypeScript", "Next.js"]

[stack.avoid]
items = ["Prisma", "ORM abstraction over raw SQL"]

[conventions]
commit_style = "conventional commits, lowercase, imperative mood"

[ai]
verbosity = "concise, skip preamble, get to the point"
tests = "write them, dont ask if I want them"
```

Your entire identity in ~300 tokens. [Full schema](docs/schema.md)

## Keeps learning

devid can watch your Claude Code sessions and pick up corrections and preferences automatically. Pre-filters for signal keywords so it only calls the API when something relevant actually happened.

```bash
export ANTHROPIC_API_KEY=sk-ant-...
devid hook install        # auto-analyze sessions when they end
devid review              # approve what it found
```

[How auto-sync works](docs/auto-sync.md)

## Weekly digest

See what your AI tools have been picking up about you:

```bash
devid digest                # last 7 days
devid digest --analyze      # suggest identity updates via API
```

[More about digests](docs/digest.md)

## Multi-machine sync

Sync your identity across machines using any git remote:

```bash
devid remote set git@github.com:you/my-devid-identity.git
devid push                  # from this machine
devid pull                  # on another machine
```

[Multi-machine setup](docs/multi-machine.md)

## All commands

Core: `init` / `distribute` / `status` / `diff` / `edit` / `doctor` / `update`

Sync: `sync` / `review` / `hook install` / `watch` / `infer` / `digest`

Multi-machine: `remote set` / `push` / `pull`

Extras: `snippet` / `mcp` / `add` / `export` / `import`

[Full command reference](docs/commands.md)

---

## Docs

- [Setup guide](docs/setup.md) - install, first run, shell completions, file locations
- [Commands](docs/commands.md) - every command and flag
- [Distribution targets](docs/targets.md) - all 13 targets and their formats
- [Identity schema](docs/schema.md) - TOML format, token budget, sensitive data handling
- [Auto-sync](docs/auto-sync.md) - session hooks, watch mode, signal filtering
- [Weekly digest](docs/digest.md) - what your AI tools learned about you
- [Multi-machine sync](docs/multi-machine.md) - push/pull across machines via git

## License

MIT

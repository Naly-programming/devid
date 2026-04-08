# Setup guide

## Install

**Linux / macOS:**

```bash
curl -fsSL https://raw.githubusercontent.com/Naly-programming/devid/main/install.sh | sh
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/Naly-programming/devid/main/install.ps1 | iex
```

**With Go (any platform):**

```bash
go install github.com/Naly-programming/devid/cmd/devid@latest
```

**Manual:** Download a binary for your platform from the [releases page](https://github.com/Naly-programming/devid/releases) and put it on your PATH.

## Create your identity

The recommended approach - let AI extract your identity from how you already work:

```bash
# With ANTHROPIC_API_KEY set (one command, no copy-paste)
export ANTHROPIC_API_KEY=sk-ant-...
devid init

# Without API key (copy-paste flow)
devid init                    # prompt copied to clipboard
                              # paste into Claude/ChatGPT, copy the TOML response
devid init --paste            # reads clipboard, saves, distributes
```

Or fill it in manually:

```bash
devid init                    # choose "Fill in manually"
```

Or import from someone else:

```bash
devid import their-identity.toml
```

## Distribute

```bash
devid distribute
```

This writes your identity to all 13 targets in one go. Run it from inside a git repo to get project-specific files.

## Verify

```bash
devid doctor                  # check everything is wired up
devid status                  # see your identity at a glance
```

## Optional: automatic sync

```bash
export ANTHROPIC_API_KEY=sk-ant-...
devid hook install            # auto-analyze sessions when they end
```

## Shell completions

```bash
# bash
devid completion bash > /etc/bash_completion.d/devid

# zsh
devid completion zsh > "${fpath[1]}/_devid"

# fish
devid completion fish > ~/.config/fish/completions/devid.fish

# powershell
devid completion powershell | Out-String | Invoke-Expression
```

## File locations

```
~/.devid/
  identity.toml        # source of truth
  queue/               # pending candidate updates
  logs/                # hook activity logs

~/.claude/
  settings.json        # hook config (after devid hook install)
  CLAUDE.md            # global identity (after devid distribute)

~/.gemini/
  GEMINI.md            # global identity (after devid distribute)
```

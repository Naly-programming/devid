# Multi-machine sync

devid can sync your identity across machines using a git remote as the backend. Your `~/.devid/` directory becomes a git repo that pushes and pulls to any git host.

## Setup

Create a private repo on GitHub (or GitLab, or wherever):

```bash
gh repo create my-devid-identity --private
```

Then configure devid to use it:

```bash
devid remote set git@github.com:you/my-devid-identity.git
devid push
```

## On another machine

```bash
devid remote set git@github.com:you/my-devid-identity.git
devid pull        # fetches identity.toml and redistributes
```

## Day to day

```bash
devid push        # after making changes
devid pull        # on another machine to pick them up
```

`devid pull` automatically runs `devid distribute` after fetching, so your targets are updated immediately.

## What syncs

- `identity.toml` - your identity
- `.gitignore` - excludes queue, logs, and temp files

What does not sync (stays local per machine):

- `queue/` - pending review candidates
- `logs/` - hook activity
- `.last_scan` - watch timestamp

## How it works

`devid remote set` initialises `~/.devid/` as a git repo with a `.gitignore` that excludes ephemeral data. `devid push` stages, commits, and pushes. `devid pull` fetches, rebases, and redistributes.

Conflicts are handled by git's normal merge strategy. Since identity.toml is a simple TOML file, conflicts are rare and easy to resolve.

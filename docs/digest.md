# Weekly digest

devid can summarise what your AI tools have learned about you over a period of time. It scans Claude Code session transcripts, extracts preference signals, and shows you a report.

## Usage

```bash
devid digest                # last 7 days
devid digest --days 30      # last 30 days
devid digest --analyze      # also suggest identity updates via API
```

## What it shows

- Total sessions and messages scanned
- Number of preference signals detected
- Signals broken down by project
- The actual messages that triggered signal detection
- Sessions ordered by activity

## With --analyze

When `ANTHROPIC_API_KEY` is set and `--analyze` is passed, devid sends the signal messages to the API with your current identity for context. The API identifies anything novel that isn't already captured in your identity.toml and queues it for review.

```bash
devid digest --analyze
devid review                # approve or reject suggested changes
```

## What counts as a signal

The same keywords used by the session-end hook:

- Corrections: "don't", "stop", "no", "not like that"
- Preferences: "prefer", "always", "instead", "rather"
- Style: "be more", "be less", "from now on"

Sessions with zero signals are counted but not analyzed. The digest is a good way to see whether the automatic sync is catching everything or if you need to run a manual sync.

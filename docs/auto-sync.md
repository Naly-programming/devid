# Automatic session sync

devid can analyze your Claude Code sessions for corrections and preference changes, queuing them for review without any manual effort.

Requires `ANTHROPIC_API_KEY` to be set.

## Session-end hook

Fires automatically when a Claude Code session closes.

```bash
export ANTHROPIC_API_KEY=sk-ant-...
devid hook install
```

This adds a `SessionEnd` hook to `~/.claude/settings.json`. When a session ends, devid reads the transcript, filters for signals, and queues any novel preferences.

Check what the hook has been doing:

```bash
devid hook logs
```

## Watch mode

Scans all recent session transcripts. Alternative to the hook, useful if you want to run it on a schedule.

```bash
devid watch --once            # scan once and exit (good for cron)
devid watch --interval 300    # continuous, scan every 5 minutes
```

## How signal filtering works

devid does not send your entire session to the API. It pre-filters for high-signal keywords in your messages:

- Corrections: "don't", "stop", "no", "not like that", "wrong"
- Preferences: "prefer", "always", "instead", "rather"
- Style: "be more", "be less", "from now on", "going forward"

If no signals are found, no API call is made. Zero tokens spent on sessions where nothing identity-relevant happened.

When signals are found, only the matching messages and their surrounding context are sent with a focused diff prompt. The API returns only new or changed fields - not a full re-extraction.

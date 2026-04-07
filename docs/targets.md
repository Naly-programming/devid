# Distribution targets

devid writes your identity to all of these from a single TOML file.

## File-based targets

These are written automatically by `devid distribute` when you're inside a git repo.

| Tool | Path | Format |
|------|------|--------|
| Claude Code (global) | `~/.claude/CLAUDE.md` | Markdown with section markers |
| Claude Code (project) | `{repo}/CLAUDE.md` | Markdown with section markers |
| Gemini (global) | `~/.gemini/GEMINI.md` | Markdown with section markers |
| Gemini (project) | `{repo}/GEMINI.md` | Markdown with section markers |
| GitHub Copilot | `{repo}/.github/copilot-instructions.md` | Markdown with section markers |
| Cursor | `{repo}/.cursor/rules/devid.mdc` | MDC with YAML frontmatter |
| Cline | `{repo}/.clinerules` | Markdown with section markers |
| Roo Code | `{repo}/.roo/rules/devid.md` | Markdown (devid owns file) |
| Windsurf | `{repo}/.windsurf/rules/devid.md` | Markdown with YAML frontmatter |
| Aider | `{repo}/CONVENTIONS.md` | Markdown with section markers |
| AGENTS.md | `{repo}/AGENTS.md` | Markdown with section markers |

## Non-file targets

| Tool | Method |
|------|--------|
| Claude.ai | `devid snippet` copies to clipboard, or `devid mcp` for direct access |
| ChatGPT | `devid snippet` for custom instructions, `devid snippet --json` for API |
| Any MCP client | `devid mcp` starts a JSON-RPC server on stdin/stdout |

## Section markers

Files that use section markers wrap devid content between:

```
<!-- devid:start -->
<!-- managed by devid - do not edit between markers -->
...your identity...
<!-- devid:end -->
```

Anything you write outside these markers is preserved across updates. Cursor, Roo Code, and Windsurf use standalone files that devid owns entirely (no markers needed).

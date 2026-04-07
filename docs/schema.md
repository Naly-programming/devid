# Identity schema

The TOML schema is designed for maximum signal per token. Values are fragments, not sentences.

```toml
[meta]
version = "1"

[identity]
name = "Nathan"
tone = "direct, plain-spoken, no fluff, northern"
comments = "sound like the dev wrote it, not a textbook"
responses = "prose over bullets, hyphens not em dashes"
pace = "move fast, skip obvious explanations"

[stack]
primary = ["Go", "TypeScript", "Next.js"]
secondary = ["C#", ".NET"]
data = ["PostgreSQL", "Supabase"]
infra = ["Docker", "GitHub Actions", "Vercel"]

[stack.avoid]
items = ["Prisma", "ORM abstraction over raw SQL"]

[stack.avoid.reasons]
playwright_go = "community wrapper, lags behind"

[conventions]
formatting = "hyphens not em dashes"
pr_style = "small focused PRs, one concern per PR"
commit_style = "conventional commits, lowercase, imperative mood"
error_handling = "explicit, no silent swallows, log with context"
naming = "clear over clever, full words not abbreviations"

[[projects]]
name = "myproject"
repo = "myproject"
stack = ["Next.js", "TypeScript", "Supabase"]
context = "B2B platform with email automation"
patterns = ["PM2 for daemon management", "Supabase auth"]

[ai]
verbosity = "concise, skip preamble, get to the point"
confirmation = "dont ask permission for obvious next steps"
suggestions = "challenge assumptions, flag alternatives"
code_comments = "minimal, only non-obvious logic"
tests = "write them, dont ask if I want them"

[learned]
entries = ["2026-04-07: prefers explicit error types over generic wrapping"]

[private]
# Never included in any distributed output
api_key = "..."
```

See `schema/identity.toml.example` for the full annotated version.

## Token budget

Target: global context under 400 tokens. Fragments average 4-6 tokens vs 15 for full sentences.

`devid distribute`, `devid init`, and `devid sync` all show token estimates after running.

## Sensitive data

The `[private]` section is loaded into memory but excluded from every rendered target. devid also scans non-private sections on save and warns if it finds patterns that look like secrets (API keys, tokens, JWTs, etc).

# Developer Identity

**Tone:** direct, plain-spoken, no fluff, no emoji, northern
**Comments:** sound like the dev wrote it, not a textbook
**Responses:** prose over bullets unless explicitly needed, hyphens not em dashes
**Pace:** move fast, minimal hand-holding, skip obvious explanations

## Stack
Go . TypeScript . Next.js . PostgreSQL . Supabase . Tailscale . PM2 . Vercel . GitHub Actions . Docker

## Avoid
- playwright-go (community wrapper, lags behind, use JS/TS for Playwright)
- Prisma
- ORM abstraction over raw SQL

## Conventions
- hyphens not em dashes, no trailing punctuation in comments
- small focused PRs, one concern per PR
- conventional commits, lowercase, imperative mood
- explicit, no silent swallows, log with context
- clear over clever, full words not abbreviations

## AI Preferences
- concise, skip preamble, get to the point
- don't ask permission for obvious next steps, just do it
- challenge assumptions, search before agreeing, flag better alternatives
- code comments: minimal, only non-obvious logic
- tests: write them, don't ask if I want them

## Learned
- 2026-04-06: prefers explicit error types in Go over generic error wrapping

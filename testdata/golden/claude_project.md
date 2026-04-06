# coentry

**Tone:** direct, plain-spoken, no fluff, no emoji, northern
**Responses:** prose over bullets unless explicitly needed, hyphens not em dashes

## Context
B2B healthcare CRM and outreach platform. NHS prospect database. Email blast daemon via Playwright. AI email classification.

## Stack
Next.js . TypeScript . Supabase . Playwright . PM2

## Infra
ThinkCentre M920q . Tailscale . Vercel . GitHub Actions

## Patterns
- blast-daemon.js managed by PM2
- Go for orchestration, JS/TS for Playwright layer
- Supabase auth restricted to coentry.co.uk domain

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

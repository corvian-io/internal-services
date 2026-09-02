# Work log

GitHub issues are the source of truth for backlog/status (`gh issue list --repo corvian-io/internal-services`) — this file is the narrative layer on top: why an item matters, what's blocking it, what to check before starting it. Don't duplicate the ticket title/state here; link the number and say something the tracker doesn't.

## Current focus

**#12 — weather: implement Weatherstack connector and mapping** (branch `feature/issue-12`, started 2026-08-11)

`GET https://api.weatherstack.com/current` → `NormalizedWeather`, plus tests for API errors and success. This is the `Fetch(ctx) ([]NormalizedX, error)` connector piece specifically — the Redis cache/shadow-key poller is separate work (#13), and the TTL/cadence calculator is separate again (#14). Before writing the mapper: confirm the actual weatherstack free-tier call limit at signup (`facts.md` flags conflicting numbers on their own site) so #14's math starts from a real number, not the ~100/month placeholder.

**Landed on this branch ahead of schedule:** `internal/shared/managers`' generic Redis manager (the cache wrapper #13's shadow-key poller will sit on top of) — not connector-shaped ingestion, so it's not really #12's scope, but it's here now. See [[decisions#redis-manager-shared]] for why it's shared rather than weather-local, and [[bugs#redis-manager-first-draft]] for what the review caught and fixed.

## Open backlog, by service

**Weather** — #12 (connector), #13 (Redis cache + shadow-key poller), #14 (TTL/cadence calculator + docs). Natural build order: #12 → #13 → #14, since #14's formula needs #13's key shapes to exist first.

**Gateway** — #1 (Section interface + registry), #2 (scatter-gather executor), #4 (search passthrough), #5 (Stocks upstream WS relay), #6 (docker-compose port exposure). #2 is the one to build carefully per [[decisions#gateway-scatter-gather]] — `WaitGroup`+mutex, not `errgroup`.

Briefing sub-feature — #27 (cache + refresh goroutine), #28 (Claude client wrapper), #29 (prompt builder), #30 (`GET /briefing` endpoint + config). Natural build order: #29 and #28 can go in parallel, #27 needs both since the refresh goroutine calls the client with a built prompt, #30 last since it just reads #27's cache.

**Jobs** — #3 (Greenhouse), #7 (Lever), #8 (Ashby), #9 (ingestion scheduler + per-company config), #10 (stale-job-closure), #11 (tech keyword filter). Three connectors (#3/#7/#8) are independent and can go in any order; #10 needs at least one connector done first to have something to reconcile against.

**Stocks** — #20 (WS reader + tradeChan), #21 (per-symbol aggregator + `Close()`), #22 (candle flush writer), #23 (watchlist DB + bootstrap endpoint), #24 (reconnect/backoff + reconcileSubscriptions). Build order matters here: #20 before #21 before #22, since each produces the input the next consumes. Watch for the per-tick-write mistake from [[bugs#stocks-per-tick-write]] when wiring #22.

**Search** — #15 (schema + GIN indexes), #16 (Streams consumer, idempotent handlers), #17 (outbox relay), #18 (`rank()` in Go), #19 (daily reconciliation sweep). #15 blocks everything else. #16/#17 are the Listen path, #19 is the independent Poll path — don't let #19 slip just because #16/#17 land first.

**Holidays** — #25 (US-only lifecycle), #26 (verify Nager API version and record the decision). Do #26 first — #25's implementation shouldn't start against an API version that might already be superseded (see [[facts#holidays]]).

## Repo gaps (not filed as tickets — flagging here so they don't get lost)

- `go.work` doesn't exist yet. Needed before any service's module can actually build against `shared/connector` locally.
- `docker-compose.yml` doesn't exist yet. Referenced by root `CLAUDE.md` as the expected local dev loop.
- `docs/project-context.md` doesn't exist yet, despite being cited as "full background" from the root `CLAUDE.md`.
- None of the five per-service checklist docs (`docs/*-service-checklist.md`) exist yet, despite being cited from every service's own `CLAUDE.md`.

## Not yet designed

music, sports, books — mentioned in the root `CLAUDE.md` tree, no checklist, no issues filed, no service directory yet.

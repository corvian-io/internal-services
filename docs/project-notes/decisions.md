# Architectural decisions

One entry per decision. Include the rejected alternative, not just the choice — that's usually the part someone re-litigates later.

---

## Repo layout: multi-module monorepo, not a single Go module

Each service under `services/*` is its own Go module with its own `go.mod`. `shared/connector` is also its own module. A root `go.work` ties them together for local dev.

**Why:** independent build/deploy per service, and it prevents accidental cross-service imports at compile time — a single shared module would make it trivially easy to reach into another service's package. Avoids per-service repo overhead, which isn't worth it for a one-person project.

**Status:** decided, not yet scaffolded — `go.work` doesn't exist in the repo yet.

---

## One Postgres instance, one schema per service

`weather.*`, `stocks.*`, `jobs.*`, etc. No service reads another service's tables directly, under any circumstances.

**Why:** enforces real data ownership between services even though they share one physical Postgres instance — the whole point of this project is practicing service boundaries, and a shared-schema shortcut would quietly undo that.

**Rejected:** separate Postgres instances per service — overkill for a one-person project's infra budget; schema-per-service gets the same isolation guarantee at the application/permissions layer.

---

## Weather is the one exception to the schema rule — Redis instead of Postgres

Weather stores its cached data in Redis, not a `weather.*` schema.

**Why:** it's cache data with a natural TTL, not durable per-service state. Forcing it into Postgres would mean modeling expiry and staleness by hand for data that doesn't need to survive a restart.

**Guardrail:** this is a documented, deliberate exception — not precedent. Any future service that wants to skip its Postgres schema needs the same kind of explicit justification, not silent drift. See [[facts#weather-redis-keys]] for the actual key shapes.

---

## Gateway is the only externally-reachable service

Every domain service — including any live/streaming endpoint — is closed to the public network. If a service wants to talk to the dashboard directly, the gateway becomes the client of that internal endpoint and relays outward itself.

**Why:** single ingress point to reason about for auth/rate-limiting/observability later, and it keeps domain services honestly internal instead of "internal in theory, dual-homed in practice."

**Corrected once already:** an earlier version of the Stocks design had it opening a WS connection directly to the dashboard for live trade relay. That broke this rule. Corrected so the gateway holds the upstream `/stocks/stream` connection itself and fans it out to dashboard clients — see [[bugs#stocks-direct-dashboard-ws]]. Any future live/streaming service gets the identical treatment.

---

## Shared connector pattern: `Fetch(ctx) ([]NormalizedX, error)`, registered not hardcoded

Every service's ingestion layer imports `shared/connector` and implements this shape per external source, added to a `Registry map[string]Fetch` rather than a switch statement.

**Why:** every connector across every service (every ATS platform inside Jobs, every sibling service inside Search) has the identical shape, so the fan-out/error-isolation logic (`FetchAll`) is written once instead of per service.

**Resolved 2026-08-12, with Go generics** — earlier drafts of this doc said "not literally generic," meaning each service would hand-copy its own `Fetch`/`Registry` types. Once actually implemented, that turned out to be a false economy: `shared/connector` now defines `Fetch[T any]` and `Registry[T any]` (struct, holding a `sync.RWMutex` + `map[string]Fetch[T]`) as real generic types, so `Register`'s incremental-add and `FetchAll`'s concurrent fan-out/error-isolation logic are written once and imported everywhere, not copy-pasted per service. `FetchAll` returns `map[string]error` keyed by source name (not the originally-sketched `[]error`), so callers can tell which source failed. `Registry` is always used by pointer (`*Registry[T]`) since it holds a mutex — see `shared/connector/CLAUDE.md` for the full shape.

**Provenance:** the overall registered-by-name-not-switch pattern originated in Corvian and is being reused as-is; the generics-based construction was worked out fresh for this repo.

---

## Services never call each other directly

Sibling-to-sibling data flow (e.g. Search pulling from Jobs) goes through the sibling's public REST API — same as the gateway would call it. Never a direct package import or cross-service DB read.

**Why:** same data-ownership motivation as the Postgres schema rule — a direct import or DB read would be an invisible coupling that bypasses the service boundary entirely.

---

## Gateway: scatter-gather via `WaitGroup` + mutex, not `errgroup` {#gateway-scatter-gather}

Every domain service is called concurrently with its own bounded timeout (2s starting point). A failed/slow service degrades to `{"status": "unavailable"}` for its section; it never fails the whole response.

**Rejected: `errgroup`.** Its default is fail-fast — the first error cancels every other in-flight goroutine, which is backwards here: one dead service would take down eight working ones. Flagged once already as the wrong tool for this shape; don't reach for it again even though it looks idiomatic.

**Pattern:** each goroutine swallows its own error into a status field instead of returning it. See `services/gateway/CLAUDE.md` for the reference implementation.

---

## Gateway: Search is not in the scatter-gather fan-out

`GET /search?q=...` routes straight through to the Search service instead of joining the Weather/Stocks/Jobs/etc. batch.

**Why:** it's a different request shape — an on-demand typed query, not one of the fixed dashboard sections.

---

## Gateway: Claude-generated daily briefing — scheduled refresh, not per-request

`GET /briefing` returns a short Claude-written narrative synthesizing the current `DashboardView` (weather, stocks, jobs, holidays). Generation happens on a background ticker inside `cmd/api` — gateway's second documented exception to "no `cmd/worker`" (see the Stocks relay decision above for the first) — not per HTTP request.

**Why scheduled + cached, not on-demand:** token cost should scale with a refresh interval, not with how often the dashboard gets loaded — same reasoning as weather's shadow-key refresh cadence. The cache is in-memory only (mutex-protected struct, no Redis), since losing it on a gateway restart is a non-issue — the next tick regenerates it.

**Rejected: generating fresh on every request.** Simplest to build, but cost scales directly with request volume — exactly what a token-cost-conscious feature shouldn't do.

**Rejected: folding the briefing into the main `DashboardView` response.** Kept as a separate endpoint instead, same treatment as `/search`. A briefing failure or staleness then never affects the core dashboard sections, and it's opt-in for callers who don't want it.

**Failure handling reuses the existing degradation pattern:** if a scheduled refresh fails (API error, rate limit, `stop_reason: "refusal"`), keep serving the last successful briefing rather than erroring — same "stale beats none" philosophy as scatter-gather's per-section `{"status": "unavailable"}` degradation. `/briefing` only reports unavailable if no generation has ever succeeded.

**Model choice: `claude-sonnet-5` at `output_config.effort: "medium"`** — an explicit, deliberate override of "default to Opus" for cost reasons, made knowingly with a paid Claude membership footing the bill. Thinking left at Sonnet 5's default (adaptive — its only on-mode).

---

## `internal/shared` gains a uniform Postgres config loader (`configs` + `constants`)

Alongside the connector pattern, `internal/shared` now also holds `configs` (env-var-driven `DbConfig` construction: host/port/user/password/dbname/sslmode/schema) and `constants` (the env var name table). Added preemptively — not yet wired into any service — for stocks, jobs, holidays, and search to share one way of building a Postgres connection config, instead of four near-identical copies.

**Why it belongs in `internal/shared` and not per-service `internal/config`:** the *shape* of "load these DB env vars into a `DbConfig`" is identical across every Postgres-backed service — same bar that justified pulling the connector pattern out. Each service still supplies its own `Schema` value and its own credentials at the env-var level, so this doesn't touch the one-schema-per-service data-ownership rule — it's shared connection-*building* code, not shared data access.

---

## Weather: weatherstack, not Open-Meteo

API is `https://api.weatherstack.com/current`, key via `access_key` query param (`WEATHERSTACK_API_KEY`).

**Corrected once already** from an earlier plan that specified Open-Meteo. If Open-Meteo shows up anywhere else in the docs, it's stale — see [[bugs#weather-open-meteo-stale]].

---

## Weather: scheduled refresh via shadow-key polling, not reactive cache-aside or keyspace notifications

A poller checks for the absence of `shadow:weather:us:<zip>` on a short interval; absence is the refresh signal. Refresh writes the real key, then the shadow key **last** — ordering matters, since "shadow key exists" must mean "a refresh actually completed."

**Rejected: Redis keyspace notifications.** A missed pub/sub event has no replay. The poll-based shadow key self-heals on the next cycle instead — a crash mid-refresh just leaves the shadow key missing.

**Same reasoning reused in Search** — see the Redis Streams decision below.

---

## Stocks: three separate concerns, kept apart on purpose

1. Live tick relay (in-memory only, never touches Postgres)
2. Candle aggregation (in-memory OHLCV struct, never touches Postgres directly)
3. Watchlist (durable, small, relational — the one genuinely boring part)

**Why:** conflating these is exactly how you end up writing to Postgres per trade tick (see [[bugs#stocks-per-tick-write]]).

---

## Stocks: two-channel split (`tradeChan` / `candleChan`)

`tradeChan` fires per trade (WS reader → processor). `candleChan` fires once per symbol per `candleDuration` (aggregator → DB writer), not per trade.

**Why:** this is the fix for the per-tick-write mistake — see [[bugs#stocks-per-tick-write]]. If a DB call is about to be added inside the trade-processing loop, it belongs in the flush path instead.

**Related guardrail:** `checkInterval` (~1s, how often the ticker looks for a finished bar) and `candleDuration` (~1min, how long a bar covers) must stay separate. Conflating them either delays flushes or turns this back into a per-second DB write.

---

## Stocks: `reconcileSubscriptions` is one function, called from three places

Finnhub does not remember subscriptions across reconnects, so subscriptions must be reconciled against the watchlist on startup, on watchlist changes, *and* on every reconnect.

**Rejected:** a separate resubscribe path specific to reconnect. One function, three call sites, instead.

---

## Jobs: daily polling only — no webhook/subscription path exists

Confirmed (not assumed): none of Greenhouse, Lever, or Ashby offer an external subscription/webhook option to a third party monitoring someone else's board — their webhooks are scoped to the hiring company's own account.

**Why it matters:** this rules out an event-driven ingestion design for Jobs specifically; daily polling isn't a stopgap here, it's the only option all three platforms support.

---

## Jobs: tech-relevance filter is a keyword heuristic, applied at ingestion

No ATS platform has a universal "tech" category — department names are whatever each company chose internally. Keyword match against title + department; only matches get normalized and stored.

**Tuning note:** if this needs adjusting, change the keyword list, not the filtering point in the pipeline.

---

## Jobs: stale-job closure via seen-set reconciliation

None of the three ATS APIs signal "this job closed" — absence from a fetch is the only signal. Each refresh cycle upserts everything fetched as open, then closes anything not seen in that cycle (`status='closed' where external_id not in seen`).

**Provenance:** this is the same pattern Corvian already validated for its own stale-job-closure needs — reused, not redesigned.

**Guardrail:** if a refresh cycle only upserts and never reconciles against what it didn't see, closed jobs accumulate forever.

---

## Holidays: hardcoded to US, no country parameter

Not config-driven, not parameterized.

**Why:** this is deliberately the simplest service in the repo. Multi-country support is a real redesign if it ever becomes a requirement, not a flag to flip later — don't over-build it now.

---

## Holidays: two lifecycles, not one poll interval

Past years are fetched once, ever, immutable. Current + next year are fetched once then re-validated on a low-frequency schedule (weekly starting point) to catch legislative date revisions.

**Rejected: Weather's shadow-key pattern.** That solves "goes stale fast, refresh constantly" — the opposite of Holidays' actual shape, which is "fetch once, almost never revisit."

---

## Holidays: no second API call for long weekends

Nager.Date has a `/LongWeekend` endpoint — not used. Once holiday dates are stored, "adjacent to a weekend" is pure day-of-week arithmetic on data already owned.

**Why:** a second external dependency isn't worth it for something computable locally from data already in hand.

---

## Search: matching lives in Postgres, ranking lives in Go — kept apart on purpose

`tsvector`/GIN handles tokenizing, stemming, and candidate lookup in Postgres (same job Solr/Elasticsearch would do). Relevance weighting (title match outranks description match) is a business decision, computed in Go's `rank()` after the query returns *which field* matched.

**Corrected once already** — an earlier design put weighting into a `generated always as (...)` column, which meant a ranking change required a migration. If a ranking change ever needs a migration again, something has drifted back into the schema — see [[bugs#search-ranking-in-schema]].

---

## Search: two independent ingestion paths (Listen + Poll) — not redundant

**Listen** (real-time, primary): each sibling commits an outbox row in the same transaction as its own write (transactional outbox pattern, reused from Commerce API). A relay per publishing service polls its outbox and publishes to Redis Stream `search-events`; Search consumes via consumer group (`XREADGROUP`/`XACK`).

**Poll** (reconciliation, backstop): a daily sweep calls each sibling's own `GET` endpoint directly, bypassing the outbox/stream entirely, and reconciles the index the same stale-closure way Jobs does.

**Why both:** Poll exists specifically to catch what Listen silently drops (a relay down for an hour, a stream trimmed before consumption). Don't remove Poll because Listen "should" be reliable — Poll is the guarantee that survives Listen not being reliable.

**Idempotency requirement:** `XREADGROUP` is at-least-once, so handlers must stay upserts keyed on `(source_service, source_id)`. If a handler ever does something non-idempotent (counter increment, list append), this guarantee breaks silently.

---

## Search: Redis Streams over Postgres `LISTEN/NOTIFY`

**Why:** same reasoning as Weather's shadow-key-over-keyspace-notifications choice — Postgres notify has no listener-absent replay; a message with nobody subscribed at that moment is gone. Redis Streams persists the log and lets a consumer group replay/reclaim. Also zero new infrastructure, since Redis is already running for Weather.

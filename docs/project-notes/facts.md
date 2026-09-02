# Project facts

Config structure, endpoints, and constants. This file is checked into git, so it holds names, shapes, and public endpoints — never secret values (API keys, passwords, connection strings with credentials). Actual secrets live in environment variables / a local `.env` that's gitignored; this file documents which variable names exist and what they're for, per each service's own `CLAUDE.md`.

## Repo scaffolding status (as of 2026-08-11)

- `go.work` — referenced by root `CLAUDE.md`, **not yet created**.
- `docker-compose.yml` — referenced by root `CLAUDE.md`, **not yet created**.
- `docs/project-context.md` — referenced by root `CLAUDE.md` as "full background", **does not exist yet**.
- Per-service checklist docs (`docs/weather-service-checklist.md`, `docs/stocks-service-checklist.md`, `docs/jobs-service-checklist.md`, `docs/holidays-service-checklist.md`, `docs/search-service-checklist.md`) — referenced from each service's `CLAUDE.md`, **none exist yet**.
- Services scaffolded with a `CLAUDE.md` only (no Go code yet): gateway, weather, stocks, jobs, holidays, search, `shared/connector`.
- Not yet designed at all: music, sports, books.

## Data ownership

- One Postgres instance, one schema per service: `weather.*`, `stocks.*`, `jobs.*`, `holidays.*`, `search.*`.
- Exception: Weather uses Redis instead of a schema (cache with TTL, not durable state).

## Weather

- Data source: weatherstack, `https://api.weatherstack.com/current`.
- Auth: `access_key` query param, value from env var `WEATHERSTACK_API_KEY`.
- Free tier: ~100 calls/month (unconfirmed — conflicting numbers seen on weatherstack's own site; confirm at signup).
- Redis keys {#weather-redis-keys}:
  - `weather:us:<zip>` — actual cached data, ~7-day safety-backstop TTL (not the freshness driver).
  - `shadow:weather:us:<zip>` — empty sentinel, ~12h TTL, the actual refresh trigger.
- Refresh cadence formula: `max_refreshes/day = monthly_quota / (locations × 30)`.
- Current numbers: 1 tracked location, ~60/100 of monthly quota held back for dev/testing → ~2 refreshes/day → 12h shadow TTL. Recompute this if the tracked-location list grows.

## Gateway

- The only service that publishes a port in `docker-compose.yml`.
- Per-service call timeout in the scatter-gather fan-out: 2s (starting point, not yet measured against real usage).
- Holds the upstream connection to Stocks' internal `/stocks/stream` WS endpoint and relays to dashboard clients itself.
- Briefing feature: package `internal/briefing` — `cache.go` (mutex-protected in-memory cache: content + `generatedAt`), `refresh.go` (ticker goroutine started in `cmd/api`'s `main()`, fires once at boot then on interval), `client.go` (Anthropic Go SDK wrapper), `prompt.go` (`DashboardView` → prompt).
- Briefing env vars: `ANTHROPIC_API_KEY`, model ID (default `claude-sonnet-5`), refresh interval (default 30 min, tunable).
- Briefing endpoint: `GET /briefing`, separate from the scatter-gather batch.

## Stocks

- Live tick data source: Finnhub (WS).
- Reconnect backoff: exponential, 1s → doubling → capped at 30s, reset to 1s on successful connect.
- Bootstrap endpoint: `GET /stocks/quotes` — bulk, whole watchlist, called once on dashboard mount only (not polled).
- Live relay endpoint: `/stocks/stream` — internal-only, reachable from the gateway only.

## Jobs

- Connectors (all public, unauthenticated, no webhook/subscription option available to a third-party monitor):
  - Greenhouse: `GET https://boards-api.greenhouse.io/v1/boards/{board_token}/jobs?content=true`
  - Lever: `GET https://api.lever.co/v0/postings/{company}?mode=json`
  - Ashby: `GET https://api.ashbyhq.com/posting-api/job-board/{company}?includeCompensation=true`
- Tracked company list is config, not discoverable via any of the three APIs.
- Ingestion cadence: daily polling (only option all three platforms support).

## Holidays

- Country: hardcoded US, no parameter.
- API — version needs confirming before implementation:
  - Documented: `https://date.nager.at/api/v3/PublicHolidays/{year}/US`
  - Possibly newer: `https://nagerholidays.com/api/v4/Holidays/US/{year}` (looks like a domain/version migration in progress as of design time).
- Refresh cadence: past years fetched once ever (immutable); current + next year re-validated weekly (starting point).

## Search

- Storage/matching: Postgres full-text search (`tsvector`/GIN), schema `search.index_entries`.
- Possible later module: Meilisearch or Typesense (deliberately deferred, not decided against).
- Event stream: Redis Stream `search-events`, consumer group via `XREADGROUP`/`XACK`.
- Reconciliation sweep cadence: daily (starting point).
- Handlers must be idempotent — upsert keyed on `(source_service, source_id)`.

## Shared connector pattern

- Location: `internal/shared/connector` (moved from a root-level `shared/connector` module — see decisions.md; verified Go's `internal/` visibility works across the sibling modules in this repo's `go.work` since it's a pure import-path-prefix check, not a same-module check).
- Interface: `type Fetch[T any] func(ctx context.Context) ([]T, error)`.
- Registry: `type Registry[T any] struct { mu sync.RWMutex; fetch map[string]Fetch[T] }`, one entry per source via `Register`, no switch statement.
- Pattern (and the transactional outbox pattern used by Search, and the stale-job-closure pattern used by Jobs) originated in Corvian and is being reused as-is.

## Shared Postgres config

- Location: `internal/shared/configs` (+ `internal/shared/constants` for the env var name table) — same `internal/shared` module as the connector pattern.
- `configs.NewConfig()` / `configs.NewDbConfig()` build a `db.DbConfig` (host/port/user/password/dbname/sslmode/schema) from env vars, panicking via `GetEnvOrPanic` on anything missing.
- Env vars: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`, `DB_SCHEMA` (see `constants.EnvKeys`).
- External dependency: `github.com/akhakpouri/gorm-kit` for the `DbConfig` type itself, plus `godotenv` for local `.env` loading.
- Not wired into any service yet — added ahead of stocks/jobs/holidays/search needing it, all four being Postgres-backed.

## Shared Redis manager

- Location: `internal/shared/managers` (same `internal/shared` module as the connector pattern and Postgres config).
- `RedisManagerI[T]` / `RedisManager[T]`: `Get`, `Save`, `Exists`, `Close` — generic over the stored value type, JSON-marshaled on `Save` and unmarshaled on `Get`.
- `NewRedisManager[T](cfg configs.RedisConfig) RedisManagerI[T]` wraps `redis/go-redis/v9`.
- Config: `configs.RedisConfig` (`Host`, `Port`, `Db`, `Password`), built via `configs.NewRedisConfig()` or as the `Redis` field on `configs.NewConfig()`'s `Config`.
- Env vars: `REDIS_HOST`, `REDIS_PORT`, `REDIS_DB` (defaults to `0` if unset), `REDIS_PASSWORD` (see `constants.EnvKeys`).
- External dependency: `github.com/redis/go-redis/v9`, declared as a direct dependency in `internal/shared/go.mod`.
- First (only, so far) consumer: Weather's shadow-key cache — see `services/weather/CLAUDE.md` for the actual key shapes and TTLs.

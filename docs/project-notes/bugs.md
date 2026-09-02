# Bug log

Runtime bugs with root cause and fix go here once services are running. Nothing has shipped yet, so the entries below are design-time corrections — mistakes caught during scoping, before code existed — carried over from the per-service `CLAUDE.md` files so the reasoning doesn't get silently re-discovered later. Format is the same either way: symptom/mistake, root cause, fix.

---

## Gateway would have used `errgroup` for scatter-gather

**Mistake:** early design reached for `errgroup` to fan out calls to domain services.

**Root cause:** `errgroup`'s default behavior is fail-fast — the first error cancels every other in-flight goroutine. For a dashboard aggregating 8 independent services, that means one dead service takes down seven working ones.

**Fix:** plain `WaitGroup` + mutex; each goroutine swallows its own error into a `SectionResult{Status: "unavailable"}` instead of returning it. See `services/gateway/CLAUDE.md` for the reference loop.

**Status:** decided, not yet implemented.

---

## Stocks: live relay had the dashboard connecting directly {#stocks-direct-dashboard-ws}

**Mistake:** an earlier version of the Stocks design had the dashboard opening a WS connection straight to Stocks' `/stocks/stream` endpoint for live trade relay.

**Root cause:** violated the "gateway is the only externally-reachable service" rule — Stocks would have been dual-homed (internal service + direct public consumer).

**Fix:** `/stocks/stream` is internal-only, reachable solely from the gateway. The gateway holds its own persistent connection to it (same connect/reconnect/backoff shape Stocks itself uses for Finnhub) and fans the relay out to dashboard clients itself. Any future live/streaming service gets identical treatment.

**Status:** decided, not yet implemented.

---

## Weather: plan referenced Open-Meteo {#weather-open-meteo-stale}

**Mistake:** an earlier plan specified Open-Meteo as Weather's data source.

**Fix:** corrected to weatherstack (`https://api.weatherstack.com/current`, `WEATHERSTACK_API_KEY`). If Open-Meteo is referenced anywhere else in the docs going forward, treat it as stale, not as a second source.

**Status:** corrected in `services/weather/CLAUDE.md`. Worth a grep across `docs/` once the checklist files exist, in case Open-Meteo survived somewhere else.

---

## Stocks: per-trade Postgres writes {#stocks-per-tick-write}

**Mistake:** an earlier version of the design saved "the latest trade" to Postgres on every incoming WS message.

**Root cause:** this is the exact per-tick-write problem the eventual channel split exists to avoid — it just moved the write behind a channel instead of removing it.

**Fix:** two channels. `tradeChan` (per trade, in-memory only: last-trade cache + live broadcast + aggregator update). `candleChan` (fires once per symbol per `candleDuration`, this is the only thing that reaches Postgres). Guardrail for later: if a DB call is being added inside the trade-processing loop, stop — it belongs in the flush path.

**Related:** an aggregator with zero trades in a period must return `nil` from `Close()`, not a zero-valued candle — a naive Go zero-value struct looks like a legitimate $0.00 candle and would get written as real data.

**Status:** decided, not yet implemented.

---

## Redis manager first draft: dependency not declared, plus a few rough edges {#redis-manager-first-draft}

**Mistake:** `internal/shared/managers/redis_manager.go` imported `github.com/redis/go-redis/v9` directly, but `internal/shared/go.mod` was never tidied for it — the dependency only showed up (as `// indirect`) in `services/weather/go.mod`, which happened to pull it in transitively.

**Root cause:** `go.work` unions the module graph across the whole workspace, so the build succeeded locally even though the wrong module declared the dependency — confirmed by building `internal/shared` standalone with `GOWORK=off`, which failed to resolve the import until `go mod tidy` was run inside `internal/shared` itself.

**Fix:** ran `go mod tidy` in `internal/shared`; it's now a direct dependency there, not an indirect one in weather's `go.mod`.

**Also caught in the same review, all fixed:**
- `Save` did an `Exists` check + `Del` before writing — dead work, since Redis `SET` already overwrites an existing key and its TTL. Removed; `Save` now marshals and `SET`s directly.
- A doc comment above `Save` read `// Connect implements [RedisManagerI]` — copy-paste leftover referencing a method that doesn't exist on this interface. Corrected to reference `Save`.
- `getHost(cfg)` duplicated `RedisConfig.GetHostAddress()`, added in the same change. Removed; `NewRedisManager` now calls `cfg.GetHostAddress()`.
- `services/weather/internal/config/config.go` imported `internal/shared/configs` twice under two different names (`configs` and `cfg`). Removed the redundant import.

**Status:** fixed, not yet committed.

---

## Search: ranking logic embedded in the schema {#search-ranking-in-schema}

**Mistake:** an earlier design put relevance weighting into a `generated always as (...)` column.

**Root cause:** conflated two different jobs — matching (mechanical text processing, legitimately Postgres's job via `tsvector`/GIN) and ranking (a business decision about which field should outrank which). Putting ranking in a generated column meant a ranking change required a migration.

**Fix:** query returns *which field* matched (title vs. description); a plain Go `rank()` function decides the score. Guardrail for later: if a ranking change ever means writing a migration again, something has drifted back into the schema.

**Status:** decided, not yet implemented.

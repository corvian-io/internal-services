# Service app layout: api/worker split, shared vs. internal

Date: 2026-08-11
Status: approved, not yet implemented

## Problem

So far every service under `services/*` is just a `CLAUDE.md` describing decisions, plus a bare `go.mod` for weather. Nothing defines how a service's actual code should be organized once it has both a REST API (called by the gateway, and by siblings per the "services never call each other directly, only via REST" rule) and a background daemon (ingestion from an external source, or from a sibling via the Search outbox/stream machinery). This spec settles that layout, plus the boundary between repo-root `shared/` and per-service internal code, before any of the six services get real code.

## Decisions

### 1. Two-tier sharing boundary

- **`shared/<pkg>`** (repo root, its own Go module) — code identical in shape across *multiple* services. Today: `shared/connector` only (`Fetch(ctx) ([]NormalizedX, error)` + `Registry`). New packages land here only when something is genuinely the same shape in more than one service, not on the basis that it looks reusable.
- **`services/<name>/internal/<pkg>`** — code shared only within one service's own two binaries. Go's compiler enforces that nothing outside `services/<name>/` can import it, which matches the repo's existing "no service reads another's internals" rule at the language level, not just by convention.

### 2. Every service gets the same binary skeleton

```
services/<name>/
├── go.mod
├── Dockerfile              # multi-stage; builds both binaries into one image
├── cmd/
│   ├── api/main.go         # thin: load config → wire internal/api → serve HTTP
│   └── worker/main.go      # thin: load config → wire internal/ingest → run loop
└── internal/
    ├── api/                # HTTP handlers/router for this service's REST surface
    ├── ingest/             # connector implementation(s) + poll/consume/sweep logic
    ├── storage/            # this service's own Redis/Postgres access
    ├── config/             # env var loading
    └── model/              # NormalizedX and other domain types
```

- `cmd/api` and `cmd/worker` both import freely from `internal/*`; neither imports the other's `cmd` package.
- `worker` is the uniform binary name across every service regardless of what its background job actually does (poll, stream-consume, scheduled sweep). The *what* is documented in that service's own `CLAUDE.md`, not encoded in the binary name — grepping for "the worker" should work the same way in every service.
- Concrete connector implementations (satisfying `shared/connector.Fetch`) live in that service's own `internal/ingest`, registered into a `Registry`, per the existing connector-pattern decision. They are not shared code even though the *interface* is.
- Gateway is the one service that gets `cmd/api` only — it owns no data and has no ingestion job, so there's nothing for a worker to do. Its persistent upstream WS connection to Stocks' `/stocks/stream` is a long-lived background goroutine started inside `cmd/api` at boot, not a separate binary — it exists to serve the relay, not to ingest into any store gateway owns.

### 3. Two documented exceptions to "api serves, worker ingests"

Both are forced by where data physically lives, not by preference — noted here so they don't get "corrected" back to the clean split later without re-deriving why they're exceptions:

- **Stocks' `cmd/worker` also serves `/stocks/stream`.** The live-trade broadcaster exists only in the ingestion process's memory. Handing it to `cmd/api` would mean adding a pub/sub hop solely to preserve "api serves all endpoints" as an absolute rule — infrastructure for a purity problem, not a real one. `cmd/api` keeps `/stocks/quotes` (bulk read from Postgres watchlist, called once per dashboard mount).
- **Search's outbox relay lives inside the *publisher's* `cmd/worker`, not inside Search.** Any sibling wired into Search's Listen path (Jobs today; Stocks/Holidays are candidates per `services/search/CLAUDE.md`, not yet wired) gets an extra job in its own worker — poll that service's outbox table, publish to the `search-events` Redis Stream — alongside its existing ingestion job. Same process, same Postgres connection already open for that service's own writes. No third binary per publishing service.

### 4. Per-service topology

| Service  | `cmd/api` serves | `cmd/worker` does |
|----------|-------------------|--------------------|
| gateway  | Scatter-gather dashboard view, `GET /search` passthrough, Stocks relay fan-out to dashboard clients | *(none — no worker binary)* |
| weather  | Reads cached `NormalizedWeather` from Redis | Shadow-key poller: weatherstack connector, writes both Redis keys |
| stocks   | `GET /stocks/quotes` (Postgres watchlist read) | Finnhub WS reader → aggregator → candle flush writer (Postgres); serves `/stocks/stream` directly (exception above); outbox relay if/when Stocks joins Search's Listen path |
| jobs     | Serves job listings; `GET` endpoint Search's Poll path reconciles against | Daily poll across Greenhouse/Lever/Ashby connectors, tech-keyword filter, stale-job closure, outbox relay publishing to `search-events` |
| holidays | Serves holiday data | Fetch-once ingestion for past years; weekly revalidation for current/next year |
| search   | `GET /search?q=...` | Redis Streams consumer (idempotent handlers) *and* daily reconciliation sweep — two jobs, one worker binary |

### 5. Docker / compose

Root `CLAUDE.md`'s existing "one Dockerfile per service" convention is unchanged — one multi-stage Dockerfile per service builds both binaries into one image. `docker-compose.yml` runs **two containers per service that has a worker** (e.g. `weather-api`, `weather-worker`), same image, differing only in `command:`. Gateway gets one container, since it has no worker. Only `cmd/api` containers — plus Stocks' `cmd/worker`, per the documented exception — actually listen on a port inside the docker network; every other worker is a listener-less long-running process. The existing "only gateway publishes a port to the host" rule is unaffected.

## Rejected alternatives

- **Single binary per service** (API + background jobs as goroutines in one process). Simpler locally (one container instead of two per service), and would have sidestepped the Stocks in-memory-coupling question entirely — but it's less operational realism to practice (independent restart/scaling/failure isolation between serving and ingestion is a real thing this project exists to practice), so it was rejected in favor of separate binaries.
- **`apps/` top-level naming instead of `cmd/`.** Reads closer to a JS/Nx-monorepo convention. Rejected — not idiomatic Go, doesn't buy anything functional, and this is a Go monorepo.
- **Per-service descriptive worker names** (`cmd/poller`, `cmd/ingest`, `cmd/consumer`, ...). More self-documenting from a directory listing alone, but breaks cross-service pattern-matching, and Search's worker does two different things anyway so one descriptive name wouldn't fit it regardless. Rejected in favor of a uniform `cmd/worker`.
- **Separate `cmd/outbox-relay` binary per publishing service.** Would let the relay's restart/latency profile move independently of that service's main ingestion job. Rejected for now as more binaries than the current data volume justifies — revisit only if a real operational reason shows up (e.g. Jobs' daily poll needing to restart without disturbing search-index freshness).

## Out of scope

- Actually scaffolding the directories/binaries/Dockerfiles for all six services — that's the implementation plan that follows this spec.
- `go.work` and `docker-compose.yml` — both still don't exist in the repo (tracked in `docs/project-notes/issues.md`); this spec assumes they'll be created as part of implementing this layout.
- Any shared HTTP router/middleware package. Every `internal/api` is independent for now; only promote something to `shared/` once a second service's `internal/api` turns out to need the identical thing, not before.

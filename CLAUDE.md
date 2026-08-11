# internal-services

Personal dashboard project aggregating data from independent microservices (weather, stocks, jobs, GitHub, holidays, search, and more to come) behind a single gateway. Built explicitly for distributed-systems practice — real service boundaries, data ownership, and an actual aggregation layer — not to ship a product. Full background: `docs/project-context.md`.

## Repo shape

Monorepo, but not a single Go module. Each service under `services/*` is its own Go module with its own `go.mod`, tied together by a root `go.work`. Shared code (the connector pattern) lives in `shared/connector` as its own module too. This gets independent build/deploy per service and prevents accidental cross-service imports, without paying per-service repo overhead for a one-person project.

```
internal-services/
├── go.work
├── docker-compose.yml
├── docs/
│   ├── project-context.md              — full architecture history and decisions
│   └── *-service-checklist.md          — per-service design docs
├── services/
│   ├── gateway/
│   ├── weather/
│   ├── stocks/
│   ├── jobs/
│   ├── github/
│   ├── holidays/
│   ├── search/
│   ├── music/       (not yet designed)
│   ├── sports/      (not yet designed)
│   └── books/       (not yet designed)
└── shared/
    └── connector/    — Fetch(ctx) interface + registry, imported by every service's ingestion layer
```

## Non-negotiable architecture rules

These apply to every service in this repo. If a change would violate one of these, stop and flag it rather than working around it — several were violated once already during design and had to be corrected (see `docs/project-context.md`).

1. **One Postgres instance, one schema per service** (`weather.*`, `stocks.*`, etc.). No service reads another service's tables directly, ever. One documented exception exists — Weather uses Redis instead of a schema (cache with a TTL, not durable state). Any new exception needs the same kind of explicit justification, not silent drift.
2. **The gateway is the only service reachable from outside the private network.** Every domain service, including any live/streaming endpoint, is closed to the public. If a service ever wants to talk to the dashboard directly, that's wrong — it goes through the gateway, which becomes the client of that internal endpoint and relays outward itself.
3. **Every service's ingestion layer uses the same connector shape**: `Fetch(ctx) ([]NormalizedX, error)`, one implementation per external source, registered rather than hardcoded. Import `shared/connector`, don't reinvent this per service.
4. **Services never call each other directly.** Sibling-to-sibling data flow (Search polling Jobs, for example) goes through the sibling's public REST API, same as the gateway would call it — never a direct package import or DB read across service boundaries.

## Conventions

- Go, multi-stage Docker builds, one Dockerfile per service under `services/<name>/`.
- `docker-compose up` is the expected local dev loop — only the gateway publishes a port to the host; every other service is reachable solely over the internal docker network.
- Config via environment variables, documented per service in that service's own `CLAUDE.md`.
- Each service has its own `CLAUDE.md` one level down — read it before working in that directory. It has the service-specific decisions this root file deliberately doesn't repeat.

## Project memory

`docs/project-notes/` tracks things that don't belong in a per-service `CLAUDE.md` because they're decisions, history, or status rather than standing rules:

- `decisions.md` — architectural decision records, including the rejected alternative for each.
- `bugs.md` — mistakes caught (in design or in code) with root cause and fix.
- `facts.md` — config structure, endpoints, and constants. Never secret values — those stay in env vars, this file documents the variable names and what they're for.
- `issues.md` — narrative work log layered on top of GitHub issues (the tracker is the source of truth for status).

Before proposing an architectural change, check `decisions.md` for whether it's already been decided (and rejected) once. When something breaks, check `bugs.md` for whether it's a repeat. When looking up a port, endpoint, or config shape, prefer `facts.md` over assuming.

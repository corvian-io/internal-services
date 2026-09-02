# services/weather

Full design doc: `docs/weather-service-checklist.md`. The one service in this repo that doesn't use Postgres — read that before assuming the repo-wide schema rule applies here unmodified.

## API: weatherstack, not Open-Meteo

`https://api.weatherstack.com/current`, API key via `access_key` query param (`WEATHERSTACK_API_KEY`), free tier ~100 calls/month (conflicting numbers seen on their own site — confirm at signup). This was corrected once already from an earlier plan that said Open-Meteo; if you see Open-Meteo referenced anywhere, it's stale.

## Storage: Redis, not a Postgres schema — on purpose

This is cache data with a natural TTL, not durable per-service state, so it's the one deliberate exception to "every service gets a `<service>.*` Postgres schema." Don't add a `weather.*` schema later without re-deriving why this exception existed in the first place.

Two keys per location:
- `weather:us:<zip>` — the actual data, ~7-day safety-backstop TTL (not the freshness driver)
- `shadow:weather:us:<zip>` — empty sentinel, ~12h TTL, the actual refresh trigger

## Ingestion: scheduled refresh via shadow key, not reactive cache-aside

Locations are a small, fixed, config-driven set — not arbitrary dashboard-supplied zips. A poller checks the shadow key's existence on a short interval (every few minutes); its absence is the refresh signal. On refresh: fetch from weatherstack, write the real key, write the shadow key **last** — that ordering matters, since "shadow key exists" is only meaningful if it means "a refresh actually completed." A crash mid-refresh just leaves the shadow key missing, which self-heals on the next poll cycle instead of faking freshness.

Chose polling over Redis keyspace notifications on purpose: a missed pub/sub event has no replay; the poll-based shadow key does.

## TTL math, if the tracked-location list ever grows

```
max_refreshes/day = monthly_quota / (locations × 30)
```
Currently 1 location, ~60/100 budget held back for dev/testing → ~2 refreshes/day → 12h shadow TTL. Recompute if more locations get added — this isn't a fixed constant, it's derived from quota.

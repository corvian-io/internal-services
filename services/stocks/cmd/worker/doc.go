// Command stocks-worker holds the Finnhub WS connection, runs the per-trade
// processing loop (tradeChan) and the per-candle flush loop (candleChan),
// and — as this service's one documented exception to "api serves, worker
// ingests" — also serves /stocks/stream directly.
//
// Reflection points:
//   - /stocks/stream lives here, not in cmd/api, because the live-trade
//     broadcaster only exists in this process's memory; moving it to api
//     would mean adding a pub/sub hop purely to preserve a rule, not to
//     solve a real problem. See
//     docs/superpowers/specs/2026-08-11-service-app-layout-design.md.
//   - Do not write to Postgres per trade. tradeChan (in-memory: last-trade
//     cache, live broadcast, aggregator update) and candleChan (fires once
//     per symbol per candleDuration, the only thing that reaches Postgres)
//     are a deliberate split — an earlier version collapsed them and wrote
//     per tick. See docs/project-notes/bugs.md.
//   - checkInterval (~1s, how often the ticker looks) and candleDuration
//     (~1min, how long a bar covers) are not the same number — conflating
//     them turns this back into a per-second DB write.
//   - reconcileSubscriptions(watchlist) runs on startup, on watchlist
//     changes, AND on every reconnect — one function, three call sites, not
//     a separate reconnect-specific path.
//   - /stocks/stream is internal-only, reachable from the gateway alone.
package main

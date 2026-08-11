// Command stocks-api serves GET /stocks/quotes: a bulk, whole-watchlist
// read from Postgres. Called exactly once, when the dashboard mounts,
// before its live connection to the gateway's relay is up.
//
// Reflection points:
//   - Not polled repeatedly — the live relay (served by cmd/worker, not
//     this binary) carries updates after the first load. If this ever gets
//     called on an interval, something upstream is wrong.
//   - /stocks/stream is NOT served here — see cmd/worker/doc.go for why.
package main

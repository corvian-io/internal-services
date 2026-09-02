// Package ingest implements the Finnhub WS reader, the per-symbol candle
// aggregator, the candle flush writer, and the live-trade broadcaster
// served by cmd/worker as /stocks/stream.
//
// Reflection points:
//   - Three separate concerns live here on purpose (live relay, candle
//     aggregation, watchlist writes) — don't conflate them back together,
//     see docs/project-notes/decisions.md.
//   - An aggregator with zero trades in a period must return nil from
//     Close(), not a zero-valued candle — a naive Go zero value looks like
//     a legitimate $0.00 candle and would get written as real data.
//   - Reconnect/backoff: exponential, 1s → doubling → capped at 30s, reset
//     to 1s on successful connect. Finnhub does not remember subscriptions
//     across reconnects.
package ingest

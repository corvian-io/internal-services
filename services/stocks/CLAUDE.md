# services/stocks

Full design doc: `docs/stocks-service-checklist.md`. The most involved service in this repo — real-time ingestion, in-memory aggregation, and a durable Postgres store all in one.

## Three separate problems, don't conflate them

1. **Live tick relay** — every trade broadcast immediately to the gateway's WS connection. In-memory only, never touches Postgres.
2. **Candle aggregation** — every trade updates an in-memory per-symbol OHLCV struct. Also never touches Postgres directly.
3. **Watchlist** — durable, small, relational. The one genuinely boring part.

## The channel split — and the mistake it corrects

Two channels: `tradeChan` (WS reader → processor, fires per trade) and `candleChan` (aggregator → DB writer, fires once per symbol per `candleDuration`, not per trade). **Do not write to Postgres per trade.** An earlier version of this design tried to save "the latest trade" to the DB on every message — that's a rewrite of the exact per-tick-write problem this split exists to avoid, just moved behind a channel instead of removed. If you're about to add a DB call inside the trade-processing loop, stop and check whether it belongs in the flush path instead.

```go
for t := range tradeChan {
    lastTrade[t.Symbol] = t         // in-memory, serves the cold-start bootstrap endpoint
    liveRelay.Broadcast(t)          // non-blocking — a slow client shouldn't stall this loop
    aggregators[t.Symbol].Update(t) // in-memory, no DB
}
```

## `checkInterval` vs. `candleDuration` — do not merge these

`checkInterval` (~1s) is how often a ticker *looks* for a bar that's done — cheap, in-memory scan, fine to run tight. `candleDuration` (~1min) is how long a bar actually *covers* — this is the number that controls write volume. Conflating the two either delays flushes unnecessarily or turns this back into a per-second DB write, which is the exact problem the split was designed to prevent.

An aggregator with zero trades in a period must return `nil` from `Close()`, not a zero-valued candle — a naive Go zero-value struct looks like a legitimate $0.00 candle and will get written as real data otherwise.

## Reconnect/backoff

Exponential, 1s → doubling → capped at 30s, reset to 1s on successful connect. Finnhub does not remember subscriptions across connections — `reconcileSubscriptions(watchlist)` has to run on every reconnect, not just on startup and on watchlist changes. It's one function called from three places; don't write a separate resubscribe path for reconnect specifically.

## Live relay is internal-only — see services/gateway/CLAUDE.md

`/stocks/stream` is reachable from the gateway, never the public. The gateway holds the upstream connection and relays to the dashboard itself.

## Bootstrap endpoint

`GET /stocks/quotes` — bulk, whole watchlist at once, not per-symbol. Called exactly once, when the dashboard mounts, before its live connection is up. Never polled repeatedly — if you find yourself calling this on an interval, something's wrong; the live relay is what carries updates after the first load.

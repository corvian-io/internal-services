// Package storage is stocks' Postgres client: watchlist reads/writes (the
// "genuinely boring" durable part) and candle writes from the flush loop.
// Never called per trade — see internal/ingest.
//
// Reflection points:
//   - Watchlist is small and relational; candles are the one thing
//     internal/ingest's flush path writes here. Live ticks never touch
//     this package at all.
package storage

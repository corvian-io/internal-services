// Package model defines stocks' normalized domain types: trade ticks,
// candles/OHLCV, and the watchlist entry shape.
//
// Reflection points:
//   - Keep the live-tick type and the durable candle type distinct even
//     though they're related — they have different consumers (broadcaster
//     vs. Postgres) and different lifetimes (transient vs. durable).
package model

// Package config loads stocks' environment variables: Finnhub credentials,
// Postgres connection info, candleDuration/checkInterval.
//
// Reflection points:
//   - candleDuration and checkInterval are separate config values, not
//     one — see internal/ingest/doc.go for why conflating them is a repeat
//     mistake waiting to happen.
package config

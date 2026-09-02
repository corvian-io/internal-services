// Package api holds the HTTP handler for GET /stocks/quotes.
//
// Reflection points:
//   - This is the only HTTP surface cmd/api serves. The other stocks
//     endpoint the gateway talks to, /stocks/stream, is served by
//     cmd/worker instead — see internal/ingest for why.
package api

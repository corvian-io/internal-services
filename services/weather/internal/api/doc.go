// Package api holds the HTTP handlers for weather's REST surface, called by
// the gateway's scatter-gather fan-out.
//
// Reflection points:
//   - Read-only: every handler here reads internal/storage, none of them
//     write. Writes belong to internal/ingest, called from cmd/worker only.
package api

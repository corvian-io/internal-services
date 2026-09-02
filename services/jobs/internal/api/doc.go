// Package api holds jobs' HTTP handlers, including the endpoint Search's
// daily reconciliation sweep calls directly (bypassing the outbox/stream
// path entirely).
//
// Reflection points:
//   - Serves both the gateway's scatter-gather and Search's Poll path — see
//     cmd/api/doc.go.
package api

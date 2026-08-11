// Command jobs-api serves job listings from Postgres, including the GET
// endpoint Search's Poll path reconciles against directly.
//
// Reflection points:
//   - Read-only. Closure of stale jobs happens in cmd/worker
//     (internal/ingest), never here.
//   - This endpoint has two different callers with two different
//     expectations: the gateway's scatter-gather (fast, dashboard-facing)
//     and Search's reconciliation sweep (bulk, backstop). Don't optimize
//     this purely for one of them.
package main

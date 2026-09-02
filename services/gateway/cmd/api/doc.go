// Command gateway-api is the only externally-reachable binary in this repo.
// It composes the dashboard view via scatter-gather, routes GET /search
// straight through to Search, and holds the upstream Stocks WS connection
// that it fans out to dashboard clients.
//
// Reflection points:
//   - Gateway has no cmd/worker — it owns no data and has no ingestion job.
//     The Stocks relay connection is a long-lived background goroutine
//     started here at boot, not a separate binary. See
//     docs/superpowers/specs/2026-08-11-service-app-layout-design.md.
//   - Only this service publishes a port to the host in docker-compose.yml.
package main

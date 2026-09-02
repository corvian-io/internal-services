// Command search-worker runs two independent jobs: the Redis Streams
// consumer (search-events, real-time primary path) and the daily
// reconciliation sweep (bypasses the stream entirely, backstop path).
//
// Reflection points:
//   - These two jobs are not redundant. The sweep exists specifically to
//     catch what the consumer silently drops (a relay down for an hour, a
//     stream trimmed before consumption) — don't remove it because the
//     consumer "should" be reliable.
//   - Consumer handlers must stay idempotent: XREADGROUP is at-least-once,
//     and handlers are upserts keyed on (source_service, source_id). If a
//     handler ever does something non-idempotent, this guarantee breaks
//     silently.
//   - The outbox relay that feeds search-events lives in each publishing
//     sibling's own worker (e.g. jobs-worker), not here — see
//     docs/superpowers/specs/2026-08-11-service-app-layout-design.md.
package main

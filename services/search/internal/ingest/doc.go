// Package ingest implements the Redis Streams consumer (Listen), the daily
// reconciliation sweep (Poll), and rank() — the Go-side relevance scoring
// that deliberately stays out of the Postgres schema.
//
// Reflection points:
//   - Chose Redis Streams over Postgres LISTEN/NOTIFY: notify has no
//     listener-absent replay, Streams persists the log and lets a consumer
//     group replay/reclaim. Same reasoning as weather's shadow-key choice.
//   - rank() should be testable with no database involved — if changing it
//     ever means writing a migration, that's a regression back to the
//     rejected design. See docs/project-notes/bugs.md.
package ingest

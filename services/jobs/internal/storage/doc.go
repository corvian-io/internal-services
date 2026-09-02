// Package storage is jobs' Postgres client: job upserts, stale-closure
// queries (status='closed' where external_id not in seen), and the outbox
// table read by the relay job.
//
// Reflection points:
//   - A refresh cycle that only upserts what it fetched and never
//     reconciles against what it didn't will accumulate closed jobs
//     forever — the reconciliation query here is not optional cleanup.
package storage

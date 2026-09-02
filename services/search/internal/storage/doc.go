// Package storage is search's Postgres client: search.index_entries,
// tsvector/GIN-backed matching only. Also the Redis Stream consumer-group
// bookkeeping (XREADGROUP/XACK) for search-events.
//
// Reflection points:
//   - Only matching lives here. If relevance logic starts leaking into a
//     query in this package, it belongs in internal/ingest's rank()
//     instead.
package storage

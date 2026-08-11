// Package ingest implements the three ATS connectors (Greenhouse, Lever,
// Ashby — one shape, three different JSON responses), the tech-keyword
// filter applied at ingestion, stale-job closure, and the outbox relay.
//
// Reflection points:
//   - Tech filter is a heuristic (keyword match against title +
//     department), applied uniformly regardless of source. If it needs
//     tuning, change the keyword list, not the filtering point in the
//     pipeline.
//   - Tracked company list is config (internal/config), not discoverable —
//     none of the three platforms expose a search/discovery endpoint.
package ingest

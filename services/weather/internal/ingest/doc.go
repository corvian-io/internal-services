// Package ingest implements weather's weatherstack connector (the
// shared/connector.Fetch shape) and the shadow-key refresh logic that
// decides when to call it.
//
// Reflection points:
//   - API is weatherstack, not Open-Meteo — an earlier plan said
//     Open-Meteo; that was corrected once already, see
//     docs/project-notes/bugs.md. If Open-Meteo shows up again anywhere,
//     it's stale.
//   - The free-tier call budget (~100/month, unconfirmed) drives the
//     shadow-key TTL via docs/project-notes/facts.md's refresh-cadence
//     formula — don't hardcode a refresh interval without re-deriving it
//     from that formula if the tracked-location list changes.
package ingest

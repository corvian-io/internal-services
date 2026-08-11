// Command weather-worker runs the shadow-key poller: on a short interval it
// checks whether shadow:weather:us:<zip> is missing, and if so calls the
// weatherstack connector and refreshes both Redis keys.
//
// Reflection points:
//   - Write order is not arbitrary: write the real key, then the shadow key
//     last. "Shadow key exists" must mean "a refresh actually completed" —
//     see docs/project-notes/decisions.md.
//   - Chose polling over Redis keyspace notifications deliberately: a
//     missed pub/sub event has no replay, a missed poll cycle self-heals on
//     the next one instead.
//   - Locations are a small, fixed, config-driven set (internal/config),
//     not arbitrary dashboard-supplied zips.
package main

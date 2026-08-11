// Package model defines search's normalized domain type, IndexEntry, and
// the event payload shapes carried over search-events (job.opened,
// job.closed, ...).
//
// Reflection points:
//   - job.opened carries enough denormalized data to build the index row
//     with no follow-up call; job.closed needs only the source ID. Keep
//     that asymmetry — don't force both events into one shape.
package model

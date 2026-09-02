// Command holidays-api serves holiday data from Postgres.
//
// Reflection points:
//   - Hardcoded to US, no country parameter — this is deliberate
//     simplicity, not a gap. See docs/project-notes/decisions.md before
//     adding one.
//   - "Is this adjacent to a weekend" is computed here from stored data,
//     not fetched from Nager.Date's /LongWeekend endpoint.
package main

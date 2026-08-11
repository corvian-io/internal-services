// Package model defines holidays' normalized domain type,
// NormalizedHoliday, covering both past (immutable) and current/next-year
// (revalidated) rows.
//
// Reflection points:
//   - One type for both lifecycles — the difference between them is
//     handling in internal/ingest, not a schema difference.
package model

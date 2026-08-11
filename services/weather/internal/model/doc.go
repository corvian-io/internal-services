// Package model defines weather's normalized domain type, NormalizedWeather
// — the shape both internal/ingest (writer) and internal/api (reader) agree
// on.
//
// Reflection points:
//   - This is the NormalizedX referenced by shared/connector's Fetch shape
//     for this service specifically — see shared/connector/CLAUDE.md.
package model

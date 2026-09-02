// Package model defines jobs' normalized domain type, NormalizedJob, shared
// across all three ATS connectors despite their different source JSON
// shapes.
//
// Reflection points:
//   - One shape in, regardless of source — normalization happens in each
//     connector inside internal/ingest, not here.
package model

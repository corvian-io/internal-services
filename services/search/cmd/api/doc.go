// Command search-api serves GET /search?q=... — an on-demand typed query,
// which is why the gateway routes it as a passthrough instead of folding it
// into the scatter-gather batch.
//
// Reflection points:
//   - Matching happens in Postgres (tsvector/GIN); ranking happens in Go.
//     If a ranking change ever needs a migration, something has drifted
//     back into the schema — see docs/project-notes/bugs.md.
package main

// Package config loads jobs' environment variables: the tracked-company
// list per ATS platform, Postgres connection info.
//
// Reflection points:
//   - The tracked-company list is config precisely because none of the
//     three ATS APIs expose a way to discover it — see internal/ingest.
package config

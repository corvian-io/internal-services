// Package config loads weather's environment variables
// (WEATHERSTACK_API_KEY, Redis connection info, the tracked-location list)
// for both cmd/api and cmd/worker.
//
// Reflection points:
//   - The tracked-location list belongs here, not in internal/ingest — it's
//     config, not something discovered or computed at runtime.
package config

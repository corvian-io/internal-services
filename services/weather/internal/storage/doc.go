// Package storage is weather's Redis client: the only package in this
// service allowed to read or write weather:us:<zip> and
// shadow:weather:us:<zip>.
//
// Reflection points:
//   - Redis, not Postgres — this service is the one documented exception to
//     the repo's one-schema-per-service rule (cache with a TTL, not durable
//     state). See docs/project-notes/decisions.md before adding a
//     weather.* schema here.
package storage

// Command weather-api serves weather's read-side REST API: it reads the
// cached NormalizedWeather value out of Redis and returns it. It never
// calls weatherstack directly and never writes to Redis — that's worker's
// job.
//
// Reflection points:
//   - No fallback to a live weatherstack call on a cache miss. A miss means
//     the shadow-key poller hasn't run yet or the real key's 7-day
//     safety-backstop TTL expired — both are worker's problem to fix, not
//     api's to paper over.
//   - Weather is the one service without a Postgres schema — see
//     docs/project-notes/decisions.md for why that exception exists before
//     assuming this package should read from anywhere else.
package main

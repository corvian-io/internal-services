// Command holidays-worker runs two independent lifecycles: past years
// (fetched once, ever, immutable) and current/next year (fetched once,
// then re-validated weekly to catch legislative date revisions).
//
// Reflection points:
//   - API version needs confirming before this is built for real: the
//     documented endpoint is date.nager.at/api/v3, but a newer
//     nagerholidays.com/api/v4 turned up during design — looks like a
//     migration in progress. See docs/project-notes/facts.md and issue #26.
//   - Do not apply weather's shadow-key pattern here — that solves "goes
//     stale fast, refresh constantly," which is the opposite of this
//     service's actual shape ("fetch once, almost never revisit").
package main

// Package model defines gateway's aggregated response shapes: the dashboard
// view and per-section result (status + data, where a failed section is
// {"status": "unavailable"} rather than an error).
//
// Reflection points:
//   - A section's failure is data (a status field), not a Go error
//     propagated up — that's what lets one dead service degrade gracefully
//     instead of failing the whole response.
package model

// Package ingest implements the Nager.Date (or successor) connector and the
// two-lifecycle fetch logic: past years once-ever, current/next year on a
// weekly revalidation.
//
// Reflection points:
//   - Confirm the API version (v3 vs. v4) before wiring the connector up —
//     don't assume the documented v3 endpoint is still current by build
//     time. See issue #26.
//   - No /LongWeekend call — that's computed in internal/api from data
//     already owned, not fetched.
package ingest

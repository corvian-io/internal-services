// Package api holds gateway's HTTP handlers: the dashboard aggregation
// endpoint and the GET /search passthrough.
//
// Reflection points:
//   - /search is not part of the scatter-gather batch — it's a different
//     request shape (on-demand typed query vs. fixed dashboard sections),
//     so it doesn't route through internal/scatter at all.
package api

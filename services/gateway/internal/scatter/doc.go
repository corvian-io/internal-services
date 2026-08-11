// Package scatter implements the scatter-gather executor: every registered
// section is called concurrently with its own bounded timeout (2s starting
// point), and a failed or slow section degrades to
// SectionResult{Status: "unavailable"} for its own part of the response
// only — it never fails the whole response.
//
// Reflection points:
//   - Do not use errgroup here. Its fail-fast default cancels every other
//     in-flight goroutine on the first error — exactly backwards when one
//     dead service shouldn't take down seven working ones. Flagged once
//     already, see docs/project-notes/bugs.md.
//   - Use WaitGroup + mutex; each goroutine swallows its own error into a
//     status field instead of returning it.
package scatter

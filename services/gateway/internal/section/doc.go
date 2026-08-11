// Package section defines the Section interface gateway's scatter-gather
// executor fans out to, the registry of sections, and the per-sibling REST
// clients backing each one (weather, stocks, jobs, holidays, ...).
//
// Reflection points:
//   - This is gateway's analog of shared/connector.Fetch — same isolation
//     principle (a failed client shouldn't block the others) — but calling
//     sibling REST APIs instead of external sources, so it's its own
//     package rather than importing shared/connector directly.
package section

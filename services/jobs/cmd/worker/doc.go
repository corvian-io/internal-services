// Command jobs-worker runs the daily poll across all three ATS connectors,
// the tech-keyword filter, stale-job closure, and — once Jobs is wired into
// Search's Listen path — the outbox relay publishing to search-events.
//
// Reflection points:
//   - Daily polling isn't a fallback pending something better: none of
//     Greenhouse, Lever, or Ashby offer a webhook/subscription option to a
//     third party monitoring someone else's board. Confirmed, not assumed —
//     see docs/project-notes/facts.md.
//   - Stale-job closure needs a seen-set reconciliation, not a bare upsert:
//     none of the three APIs signal "this job closed," absence from a fetch
//     is the only signal. Reused as-is from Corvian's own stale-job-closure
//     pattern — see docs/project-notes/decisions.md.
//   - The outbox relay job (poll this service's own outbox table, publish
//     to search-events) lives here rather than as a separate binary — see
//     docs/superpowers/specs/2026-08-11-service-app-layout-design.md.
package main

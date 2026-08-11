// Package relay holds gateway's upstream WS client to stocks' internal-only
// /stocks/stream endpoint and the non-blocking fan-out broadcaster to
// connected dashboard clients.
//
// Reflection points:
//   - The dashboard never opens a WS connection to Stocks directly — an
//     earlier version did exactly that and broke the "domain services are
//     closed to the public" rule. Gateway is the client of that internal
//     endpoint and relays outward itself. See docs/project-notes/bugs.md.
//   - Any future live/streaming service gets identical treatment: this
//     package's shape is the template, not a stocks-specific one-off.
package relay

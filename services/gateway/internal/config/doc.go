// Package config loads gateway's environment variables: per-service base
// URLs for the scatter-gather clients, the per-call timeout (2s starting
// point), and Stocks relay connection info.
//
// Reflection points:
//   - The 2s timeout is a starting point, not measured against real usage —
//     don't treat it as load-bearing without revisiting it.
package config

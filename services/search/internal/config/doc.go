// Package config loads search's environment variables: Postgres connection
// info, the Redis Stream name (search-events), and the consumer group
// name.
//
// Reflection points:
//   - No per-sibling URL list here — endpoints polled by the
//     reconciliation sweep come from each sibling's own service, not a
//     centrally duplicated config.
package config

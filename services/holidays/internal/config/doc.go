// Package config loads holidays' environment variables: Postgres
// connection info, which API base URL/version to call.
//
// Reflection points:
//   - The API version (v3/v4) is exactly the kind of thing that belongs
//     here once confirmed, so a version migration becomes a config change,
//     not a code change.
package config

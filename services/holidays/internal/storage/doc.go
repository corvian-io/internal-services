// Package storage is holidays' Postgres client.
//
// Reflection points:
//   - Past-year rows are immutable once written — nothing should ever
//     update them. Only current/next-year rows get touched by
//     revalidation.
package storage

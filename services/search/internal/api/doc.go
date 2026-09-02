// Package api holds the GET /search?q=... handler.
//
// Reflection points:
//   - This handler surfaces which field matched (title vs. description);
//     internal/ingest's rank() decides the score. Don't move relevance
//     weighting into the query itself.
package api

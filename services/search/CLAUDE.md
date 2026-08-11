# services/search

Design docs: `docs/search-service-checklist.md` for the v0/v1 storage-tech call (Postgres full-text search now, Meilisearch/Typesense as a deliberate later module). The ingestion detail below is the more complete reference for how data actually gets in — it postdates that file.

## What's actually worth indexing

Jobs is the one service with real free-text content worth ranking. Stocks' watchlist and Holidays are small, exact-match-friendly datasets — cheap to include once the machinery exists, but they're not the reason this service exists. Don't over-invest in ranking quality for tickers and holiday names.

## Matching lives in Postgres; ranking lives in Go — do not merge these back together

This was a real design correction, not a preference. `tsvector`/GIN genuinely is the search engine — tokenizing, stemming, and fast candidate lookup are legitimately offloaded to Postgres, same as Solr or Elasticsearch would do it. But relevance weighting (a title match should outrank a description match) is a business decision, not mechanical text processing, and it does not belong in a `generated always as (...)` column where it's only changeable via a migration.

```sql
create table search.index_entries (
    source_service text not null,
    source_id      text not null,
    title          text not null,
    description    text,
    url            text,
    title_vector   tsvector generated always as (to_tsvector('english', title)) stored,
    desc_vector    tsvector generated always as (to_tsvector('english', coalesce(description,''))) stored,
    primary key (source_service, source_id)
);
create index idx_title_vector on search.index_entries using gin(title_vector);
create index idx_desc_vector on search.index_entries using gin(desc_vector);
```

Query returns *which field* matched; Go decides what that means:

```go
func rank(results []Candidate) []Result {
    for i := range results {
        switch {
        case results[i].TitleHit:
            results[i].Score = 2.0
        case results[i].DescHit:
            results[i].Score = 1.0
        }
    }
    sort.Slice(results, func(i, j int) bool { return results[i].Score > results[j].Score })
    return results
}
```
If a ranking change ever means writing a migration, something has drifted back into the schema. It should always be a change to `rank()`, tested with no database involved.

## Ingestion has two independent paths — they're not redundant, they do different jobs

**Listen** (real-time, primary path): each sibling service commits an outbox row in the same transaction as its own write — this is the transactional outbox pattern from Commerce API, reused here as-is. A relay per publishing service polls its outbox and publishes to a Redis Stream (`search-events`); Search consumes via a consumer group (`XREADGROUP`/`XACK`), so an unacked message from a crashed consumer sits in the pending list instead of vanishing. `job.opened` events carry enough denormalized data to build the index row with no follow-up call; `job.closed` needs only the source ID, since deleting doesn't require knowing what the thing was.

**Poll** (reconciliation, backstop): a periodic sweep (daily is fine) calls each sibling's own `GET` endpoint directly — bypassing the outbox/stream path entirely — and reconciles the index the same stale-closure way Jobs itself reconciles against its ATS sources. This exists specifically to catch anything the listen path silently dropped: a relay down for an hour, a stream trimmed before Search consumed it. Don't remove this because listen "should" be reliable — it's the guarantee that survives listen not being reliable.

Consumer handlers must be idempotent — `XREADGROUP` is at-least-once, so the same event can be delivered twice. This works today only because handlers are upserts keyed on `(source_service, source_id)`; replaying one just writes the same row again. If a handler ever does something non-idempotent (incrementing a counter, appending to a list), this guarantee breaks silently.

## Chose Redis Streams over Postgres `LISTEN/NOTIFY`

Same reasoning as Weather's shadow-key-poll-over-keyspace-notifications choice: Postgres notify has no listener-absent replay — a message with nobody subscribed at that exact moment is gone. Redis Streams persists the log and lets a consumer group replay/reclaim, which is the actual requirement here. Also zero new infrastructure — Redis is already running for Weather.

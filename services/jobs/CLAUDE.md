# services/jobs

Full design doc: `docs/jobs-service-checklist.md`.

## Three connectors, one shape, three different JSON responses

- Greenhouse: `GET https://boards-api.greenhouse.io/v1/boards/{board_token}/jobs?content=true`
- Lever: `GET https://api.lever.co/v0/postings/{company}?mode=json`
- Ashby: `GET https://api.ashbyhq.com/posting-api/job-board/{company}?includeCompensation=true`

All public, unauthenticated, no discovery/search endpoint on any of them — the tracked company list is config, not something the service can find on its own. Confirmed (not assumed): none of the three offer an external subscription/webhook option to a third party monitoring someone else's board — their webhooks are scoped to the hiring company's own account. Daily polling is the only option here, on all three platforms, not a fallback pending something better.

## Tech filter is a heuristic, applied at ingestion

No platform has a universal "tech" category — department names are whatever each company chose internally. Keyword match against title + department, applied uniformly regardless of source. Only matching jobs get normalized and stored at all. If this needs tuning, the keyword list is the thing to adjust, not the filtering point in the pipeline.

## Stale-job-closure — this is the part that's easy to get wrong

None of the three APIs signal "this job closed." A filled or pulled posting just silently stops appearing in the feed. Absence is the only closure signal that exists:

```go
seen := make(map[string]bool)
for _, j := range fetchedJobs {
    seen[j.ExternalID] = true
    db.UpsertJob(j, "open")
}
db.CloseUnseenJobs(company, seen) // status='closed' where external_id not in seen
```

If a refresh cycle only upserts what it fetched and never reconciles against what it didn't, closed jobs accumulate in the index forever. This is the same pattern Corvian already validated for its own stale-job-closure needs — reuse it, don't redesign it.

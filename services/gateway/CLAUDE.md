# services/gateway

The only service in this repo reachable from outside the private network. Composes the dashboard view; holds no data of its own. Full design conversation in `docs/project-context.md`.

## Scatter-gather, not sequential

Calls every domain service concurrently, each with its own bounded timeout (2s is the starting point, not measured against real usage). A failed or slow service degrades to `{"status": "unavailable"}` for that section — it never fails the whole response.

**Do not use `errgroup` for this.** Its default behavior is fail-fast: the first error cancels every other in-flight goroutine, which is exactly backwards here — one dead service would take down eight working ones with it. Use a plain `WaitGroup` + mutex, and have each goroutine swallow its own error into a status field instead of returning it. This was flagged once already as the wrong tool for the job; don't reach for `errgroup` here even though it looks like the idiomatic choice for this shape.

```go
for _, s := range registeredSections {
    wg.Add(1)
    go func(s Section) {
        defer wg.Done()
        callCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
        defer cancel()
        data, err := s.Fetch(callCtx)
        mu.Lock()
        defer mu.Unlock()
        if err != nil {
            view.Sections[s.Name] = SectionResult{Status: "unavailable"}
            return
        }
        view.Sections[s.Name] = SectionResult{Status: "ok", Data: data}
    }(s)
}
```

## Search is not in that fan-out

`GET /search?q=...` routes straight through to the Search service. It's a different request shape — an on-demand typed query, not one of the fixed dashboard sections — so it doesn't belong in the scatter-gather batch alongside Weather/Stocks/Jobs/etc.

## Stocks' live relay — gateway is the client, not the dashboard

Stocks broadcasts live trades over an internal-only `/stocks/stream` WS endpoint, reachable only from the gateway. The gateway holds its own persistent connection to that endpoint — same connect/reconnect/backoff shape Stocks itself uses for its Finnhub connection — and fans the relay out to connected dashboard clients itself. The dashboard never opens a WS connection to Stocks directly. This was a real correction, not the original design — an earlier version had Stocks talking to the dashboard directly, which broke the "domain services are closed to the public" rule. If a future service ever wants a live/streaming path to the dashboard, it gets the identical treatment: gateway holds the upstream connection, dashboard only ever talks to the gateway.

## Deployment

Only this service publishes a port in `docker-compose.yml`. Every domain service (weather, stocks, jobs, github, holidays, search, and whatever comes after) is reachable solely over the internal docker network, by this service — never by the host, never by the dashboard UI directly.

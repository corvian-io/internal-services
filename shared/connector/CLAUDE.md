# shared/connector

The one piece of code every service's ingestion layer imports. Do not duplicate this pattern locally in a service instead of importing it — the whole point of pulling it out is that every connector across every service (every ATS platform inside Jobs, every sibling service inside Search) has the identical shape.

## The interface

```go
type Fetch func(ctx context.Context) ([]NormalizedX, error)
```

Not literally generic — each service defines its own `NormalizedX` (`NormalizedWeather`, `NormalizedJob`, `IndexEntry`, etc.) and its own concrete `Fetch` functions matching that shape. What's shared is the *pattern* and the registry helpers, not a single generic type.

## Registry pattern

One file per source, registered by name, not a giant switch statement:

```go
type Registry map[string]Fetch

func (r Registry) FetchAll(ctx context.Context) ([]NormalizedX, []error) {
    // fan out, collect results and errors per source — a failed source
    // shouldn't block the others, same isolation principle as the gateway's scatter-gather
}
```

## What varies per service, deliberately

- **Trigger**: most connectors are timer-driven (Weather's shadow-key refresh, Jobs' daily poll), but Search's ingestion is triggered both by a Redis Stream consumer (event-driven) and a periodic reconciliation sweep — same `Fetch` shape, different caller. Don't assume every connector is called from a ticker.
- **Source**: usually an external API, but Search's connectors point at sibling services' own REST endpoints instead — same interface, internal target.

## Validated in

This pattern originated in Corvian and is being reused here as-is, not redesigned. If something about it feels wrong for a specific service, that's more likely a sign the service doesn't fit the pattern (worth flagging) than a sign the pattern needs changing.

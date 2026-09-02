# shared/connector

The one piece of code every service's ingestion layer imports. Do not duplicate this pattern locally in a service instead of importing it — the whole point of pulling it out is that every connector across every service (every ATS platform inside Jobs, every sibling service inside Search) has the identical shape.

## The interface

```go
type Fetch[T any] func(ctx context.Context) ([]T, error)
```

This is a real Go generic type — imported, not hand-copied. Each service instantiates it with its own concrete normalized type (`Fetch[NormalizedWeather]`, `Fetch[NormalizedJob]`, `Fetch[IndexEntry]`), but `Fetch` itself and the registry machinery below are genuinely shared code. An earlier version of this doc said "not literally generic" — that described the design before implementation; it was built with generics on purpose once it came time to write it, so every service gets the fan-out/error-isolation logic for free instead of re-implementing it per service.

## Registry pattern

One file per source, registered by name via `Register`, not a giant switch statement:

```go
type Registry[T any] struct {
    mu    sync.RWMutex
    fetch map[string]Fetch[T]
}

func NewRegistry[T any]() *Registry[T]
func (r *Registry[T]) Register(name string, f Fetch[T])
func (r *Registry[T]) FetchAll(ctx context.Context) ([]T, map[string]error)
```

`FetchAll` fans out to every registered `Fetch` concurrently — `WaitGroup` + a plain `sync.Mutex` guarding the shared results, no `errgroup`, same reasoning as the gateway's scatter-gather: one failing source shouldn't cancel the others. It returns the combined results plus a `map[string]error` keyed by source name, so a caller can tell *which* source failed, not just that something did.

`Registry` holds a `sync.RWMutex` and is always used by pointer (`*Registry[T]`), never copied — `NewRegistry` returns a pointer for exactly this reason. Copying a `Registry[T]` value would silently duplicate the mutex, and mutations through one copy would stop being visible to the other.

## What varies per service, deliberately

- **Trigger**: most connectors are timer-driven (Weather's shadow-key refresh, Jobs' daily poll), but Search's ingestion is triggered both by a Redis Stream consumer (event-driven) and a periodic reconciliation sweep — same `Fetch` shape, different caller. Don't assume every connector is called from a ticker.
- **Source**: usually an external API, but Search's connectors point at sibling services' own REST endpoints instead — same interface, internal target.

## Validated in

This pattern originated in Corvian and is being reused here as-is, not redesigned. If something about it feels wrong for a specific service, that's more likely a sign the service doesn't fit the pattern (worth flagging) than a sign the pattern needs changing.

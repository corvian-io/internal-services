package connector

import (
	"context"
	"sync"
)

type Fetch[T any] func(ctx context.Context) ([]T, error)

// the register holds on to a map of key/fetch methods
type Registry[T any] struct {
	mu    sync.RWMutex
	fetch map[string]Fetch[T]
}

func NewRegistry[T any]() *Registry[T] {
	return &Registry[T]{}
}

// inserting f into fetch; a list of Fetch[T]. This is a placeholder for all fetch calls.
func (r *Registry[T]) Register(name string, f Fetch[T]) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.fetch == nil {
		r.fetch = make(map[string]Fetch[T])
	}
	r.fetch[name] = f
}

// for each name, fetch pair in r.fetch, go swpans a goroutine that calls fetch(ctx) one-by-one
func (r *Registry[T]) FetchAll(ctx context.Context) ([]T, map[string]error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	type event struct {
		key string
		fn  Fetch[T]
	}

	snapshot := make([]event, 0, len(r.fetch))
	for key, f := range r.fetch {
		snapshot = append(snapshot, event{
			key: key,
			fn:  f,
		})
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	var result []T
	errors := make(map[string]error)

	for _, e := range snapshot {
		wg.Add(1)
		go func(key string, f Fetch[T]) {
			defer wg.Done()
			entry, err := f(ctx)
			mutex.Lock()
			defer mutex.Unlock()
			if err != nil {
				errors[key] = err
			} else {
				result = append(result, entry...)
			}
		}(e.key, e.fn)
	}
	wg.Wait()
	return result, errors
}

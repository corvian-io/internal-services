package connector

import (
	"context"
	"sync"
)

type Fetch[T any] func(ctx context.Context) ([]T, error)

type Registry[T any] struct {
	mu sync.RWMutex
	fetch map[string]Fetch[T]
}

func NewRegistry[T any](f map[string]Fetch[T]) Registry[T] {
	return Registry[T]{
		fetch: f,
	}
}



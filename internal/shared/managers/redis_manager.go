package managers

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/corvian/argus/internal/shared/configs"
	"github.com/redis/go-redis/v9"
)

type RedisManagerI[T any] interface {
	Get(ctx context.Context, key string) (T, error)
	Save(ctx context.Context, key string, value T, ttl time.Duration) error
	Exists(ctx context.Context, key string) (bool, error)
	Close() error
}

type RedisManager[T any] struct {
	client *redis.Client
}

func (r *RedisManager[T]) Get(ctx context.Context, key string) (T, error) {
	var value T

	result, err := r.client.Get(ctx, key).Result()
	if err != nil {
		return value, err
	}

	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		return value, err
	}

	return value, nil
}

func (r *RedisManager[T]) Close() error {
	return r.client.Close()
}

// Save implements [RedisManagerI].
func (r *RedisManager[T]) Save(ctx context.Context, key string, value T, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		slog.Error("exception occured when creating the object", "error", err)
	}
	return err
}

// Exists implements [RedisManagerI].
func (r *RedisManager[T]) Exists(ctx context.Context, key string) (bool, error) {
	result, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func NewRedisManager[T any](cfg configs.RedisConfig) RedisManagerI[T] {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.GetHostAddress(),
		Password: cfg.Password,
		DB:       cfg.Db,
	})
	return &RedisManager[T]{
		client: client,
	}
}

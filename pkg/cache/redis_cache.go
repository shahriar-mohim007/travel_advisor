package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// NewRedis return a new redis cache
func NewRedis(client *redis.Client) Cache {
	return &Redis{client: client}
}

// Redis represents a concrete redis
type Redis struct {
	client *redis.Client
}

// Ping ping the redis redis if success return nil
func (r *Redis) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}

// Set set a key in redis
func (r *Redis) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return r.client.Set(ctx, key, value, exp).Err()
}

// Get get a key from redis
func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	resStr, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("key %s not found", key)
		}
		return "", err
	}
	return resStr, nil
}

func (r *Redis) Keys(ctx context.Context) ([]string, error) {
	keys, err := r.client.Keys(ctx, "*").Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

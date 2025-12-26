package cache

import (
	"context"
	"time"
)

type Cache interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Keys(ctx context.Context) ([]string, error)
}

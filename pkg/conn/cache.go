package conn

import (
	"context"
	"time"

	"travel_advisor/pkg/cache"
	"travel_advisor/pkg/config"
	"travel_advisor/pkg/log"

	"github.com/redis/go-redis/v9"
)

var defaultCache cache.Cache
var redisClient *redis.Client

// GetRedis return defautl connected redis client
func GetRedis() *redis.Client {
	return redisClient
}

// ConnectCache ...
func ConnectCache(cfg *config.RedisConfig) error {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	})
	defaultCache = cache.NewRedis(rdb)
	redisClient = rdb
	return rdb.Ping(ctx).Err()
}

// ConnectDefaultCache connect with default configurations
func ConnectDefaultCache() error {
	cfg := config.Redis()
	err := ConnectCache(cfg)
	// run a background process to ping and establish connection
	go func() {
		for {
			if err := defaultCache.Ping(context.Background()); err != nil {
				log.Warn("cache: ping error:", err)
				if err := ConnectCache(cfg); err != nil {
					log.Warn("cache:failed to reconnect:", err)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}()
	return err
}

// DefaultCache return default connected cache
func DefaultCache() cache.Cache {
	return defaultCache
}

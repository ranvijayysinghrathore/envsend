package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/ranvijayysinghrathore/envsend/backend/internal/config"
)

// RedisClient handles Redis operations for rate limiting and caching.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client.
func NewRedisClient(cfg config.RedisConfig) (*RedisClient, error) {
	opt, err := redis.ParseURL(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	if cfg.Password != "" {
		opt.Password = cfg.Password
	}

	// Production Tuning: Connection Pool Settings
	opt.PoolSize = 20
	opt.MinIdleConns = 10
	opt.ConnMaxIdleTime = 5 * time.Minute

	client := redis.NewClient(opt)

	// Test connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis: %w", err)
	}

	return &RedisClient{client: client}, nil
}

// Close closes the Redis connection.
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// CheckRateLimit checks if an IP address has exceeded the rate limit.
// Returns true if allowed, false if rate limited.
func (r *RedisClient) CheckRateLimit(ctx context.Context, ip string, limit int, window time.Duration) (bool, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)

	// Increment counter
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to increment rate limit: %w", err)
	}

	// Set expiry on first request
	if count == 1 {
		if err := r.client.Expire(ctx, key, window).Err(); err != nil {
			return false, fmt.Errorf("failed to set expiry: %w", err)
		}
	}

	// Check if limit exceeded
	return count <= int64(limit), nil
}

// GetRateLimitRemaining returns the number of requests remaining for an IP.
func (r *RedisClient) GetRateLimitRemaining(ctx context.Context, ip string, limit int) (int, error) {
	key := fmt.Sprintf("ratelimit:%s", ip)

	count, err := r.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return limit, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get rate limit: %w", err)
	}

	remaining := limit - count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, nil
}

// SetCache sets a value in the cache with expiry.
func (r *RedisClient) SetCache(ctx context.Context, key string, value interface{}, expiry time.Duration) error {
	return r.client.Set(ctx, key, value, expiry).Err()
}

// GetCache gets a value from the cache.
func (r *RedisClient) GetCache(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	}
	return val, err
}

// DeleteCache deletes a key from the cache.
func (r *RedisClient) DeleteCache(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// AcquireLock acquires a distributed lock.
func (r *RedisClient) AcquireLock(ctx context.Context, lockKey string, ttl time.Duration) (bool, error) {
	return r.client.SetNX(ctx, lockKey, "locked", ttl).Result()
}

// ReleaseLock releases a distributed lock.
func (r *RedisClient) ReleaseLock(ctx context.Context, lockKey string) error {
	return r.client.Del(ctx, lockKey).Err()
}

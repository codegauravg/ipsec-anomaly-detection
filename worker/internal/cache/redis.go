package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(redisURL string) *Cache {
	// In production, parse the URL properly or use Options
	client := redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
	return &Cache{client: client}
}

// IncrementAndExpire increments the counter for an IP and resets its expiration.
// Returns the new count.
func (c *Cache) IncrementAndExpire(ctx context.Context, key string, window time.Duration) (int64, error) {
	pipe := c.client.Pipeline()

	incr := pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Val(), nil
}

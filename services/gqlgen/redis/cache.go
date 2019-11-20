package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

// Cache struct
type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

const apqPrefix = "apq:"
const userPrefix = "user:"

// NewCache func
func NewCache(redisAddress string, password string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &Cache{client: client, ttl: ttl}, nil
}

// Add func
func (c *Cache) Add(ctx context.Context, hash string, query string) {
	c.client.Set(apqPrefix+hash, query, c.ttl)
}

// Get func
func (c *Cache) Get(ctx context.Context, hash string) (string, bool) {
	s, err := c.client.Get(apqPrefix + hash).Result()
	if err != nil {
		return "", false
	}
	return s, true
}

// AddIP func
func (c *Cache) AddIP(userID string, userIP string) {
	c.client.Set(userPrefix+userID, userIP, c.ttl)
}

// GetIP returns addresses IP
func (c *Cache) GetIP(userID string) (string, bool) {
	s, err := c.client.Get(userPrefix + userID).Result()
	if err != nil {
		return "", false
	}
	return s, true
}

// DeleteIP func
func (c *Cache) DeleteIP(userID string) {
	_ = c.client.Del(userPrefix + userID)
}

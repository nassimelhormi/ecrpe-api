package redis

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis"
)

type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

const apqPrefix = "apq:"

func NewCache(redisAddress string, password string, ttl time.Duration) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})
	if err := client.Ping().Err(); err != nil {
		return nil, errors.WithStack(err)
	}
	return &Cache{client: client, ttl: ttl}, nil
}

func (c *Cache) Add(ctx context.Context, hash string, query string) {
	c.client.Set(apqPrefix+hash, query, c.ttl)
}

func (c *Cache) Get(ctx context.Context, hash string) (string, bool) {
	s, err := c.client.Get(apqPrefix + hash).Result()
	if err != nil {
		return "", false
	}
	return s, true
}

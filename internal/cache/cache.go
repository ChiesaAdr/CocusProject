package cache

import (
	"cocus/internal/config"

	"github.com/go-redis/redis"
)

type Cache struct {
	order uint64
	redis *redis.Client
	limit int
}

func NewCache() *Cache {
	//History(Cache) is limited to 20 messages by default
	limit := config.GetRedisLimit()
	cache := Cache{order: 0, limit: limit}
	cache.NewRedisClient()

	return &cache
}

//Using the order ID to sort the cache
func (c *Cache) IncrOrder() uint64 {
	c.order++
	return c.order
}

//Using to unit test
func (c *Cache) GetRedis() *redis.Client {
	return c.redis
}

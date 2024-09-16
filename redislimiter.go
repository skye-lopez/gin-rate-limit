package limiter

import (
	"github.com/redis/go-redis/v9"
)

type RedisLimiter struct {
	Client *redis.Client
}

// Creates a new redis based limiter.
func NewRedisLimiter(client *redis.Client) *RedisLimiter {
	return &RedisLimiter{
		Client: client,
	}
}

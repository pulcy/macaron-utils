package redis

import (
	ms "github.com/go-macaron/session"
	redis "github.com/go-macaron/session/redis"
	mu "github.com/pulcy/macaron-utils/session"
)

func init() {
	ms.Register("retry-redis", mu.NewRetryProvider(&redis.RedisProvider{}))
}

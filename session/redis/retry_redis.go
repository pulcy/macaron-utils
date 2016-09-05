package redis

import (
	ms "github.com/go-macaron/session"
	"github.com/go-macaron/session/redis"
	"github.com/pulcy/macaron-utils/session"
)

func init() {
	ms.Register("retry-redis", session.NewRetryProvider(&redis.RedisProvider{}))
}

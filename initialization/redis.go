package initialization

import (
	"fmt"
	"github.com/go-redis/redis"
	"hertz_demo/conf"
)

type Redis struct {
	*redis.Client
}

func InitRedis(c *conf.RedisConfig) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: c.Password, // no password set
		DB:       c.DB,       // use default DB
		PoolSize: c.PoolSize,
	})
	return &Redis{
		client,
	}
}

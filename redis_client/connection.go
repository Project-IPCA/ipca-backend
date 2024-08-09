package redis_client

import (
	"fmt"

	"github.com/Project-IPCA/ipca-backend/config"
	"github.com/redis/go-redis/v9"
)

func RedisClient(cfg *config.Config) *redis.Client {
	url := fmt.Sprintf("redis://%s:%s@%s:%s/",cfg.REDIS.User,cfg.REDIS.Password,cfg.REDIS.Host,cfg.REDIS.Port)
	
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}

	return redis.NewClient(opt)
}

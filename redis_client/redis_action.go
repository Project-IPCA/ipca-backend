package redis_client

import (
	"fmt"
	"context"
	"time"
	"github.com/redis/go-redis/v9"
)

type IRedisAction interface {
	PublishMessage(channel, message string)
	SubscribeTopic(channel string)
}

type RedisAction struct {
	Redis *redis.Client
}

func NewRedisAction(redis *redis.Client) *RedisAction{
	return &RedisAction{Redis: redis}
}

func (redisAction *RedisAction )PublishMessage(channel, message string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	
    defer cancel()  
    err := redisAction.Redis.Publish(ctx, channel, message).Err()
    if err != nil {
        return fmt.Errorf("failed to publish message: %w", err)
    }

    return nil
}


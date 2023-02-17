package dao

import (
	"context"

	"github.com/go-redis/redis/v8"
)

// InitRedis 初始化redis
func InitRedis() error {
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

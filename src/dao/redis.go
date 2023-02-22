package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

// RedisClient Redis 连接
var RedisClient *redis.Client

// InitRedis 初始化redis
func InitRedis() error {
	conf := GetConfig().Redis

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", conf.Url, conf.Port),
		Password: conf.Password,
		DB:       conf.Db,
	})

	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	return nil
}

// CloseRedis 关闭redis连接
func CloseRedis() {
	err := RedisClient.Close()
	if err != nil {
		return
	}
}

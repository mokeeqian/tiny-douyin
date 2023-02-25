package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RedisClient *redis.Client

//var RedisUserRelationClient *redis.Client
//var RedisFavoriteClient *redis.Client
//var RedisCommentClient *redis.Client

// InitRedis 初始化redis
func InitRedis() {
	conf := GetConfig().Redis
	addr := fmt.Sprintf("%s:%s", conf.Url, conf.Port)

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       conf.Db,
	})

	// 使用多个db，会导致数据一致性问题
	//RedisUserRelationClient = redis.NewClient(&redis.Options{
	//	Addr:     addr,
	//	Password: conf.Password,
	//	DB:       1,
	//})
	//
	//RedisFavoriteClient = redis.NewClient(&redis.Options{
	//	Addr:     addr,
	//	Password: conf.Password,
	//	DB:       2,
	//})
	//
	//RedisCommentClient = redis.NewClient(&redis.Options{
	//	Addr:     addr,
	//	Password: conf.Password,
	//	DB:       3,
	//})
}

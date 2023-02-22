package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RedisFollowers *redis.Client
var RedisFollowing *redis.Client
var RdbFollowingPart *redis.Client // 用户关注粉丝

var RedisFavoriteUser2Video *redis.Client //key:userId,value:VideoId
var RedisFavoriteVideo2User *redis.Client //key:VideoId,value:userId

var RedisVideo2Comment *redis.Client //redis db11 -- video_id + comment_id
var RedisComment2Video *redis.Client //redis db12 -- comment_id + video_id

// InitRedis 初始化redis
func InitRedis() {
	conf := GetConfig().Redis
	addr := fmt.Sprintf("%s:%s", conf.Url, conf.Port)

	RedisFollowers = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       0,
	})
	RedisFollowing = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       1,
	})
	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       2,
	})

	RedisFavoriteUser2Video = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       3,
	})

	RedisFavoriteVideo2User = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       4,
	})
	RedisVideo2Comment = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       5,
	})

	RedisComment2Video = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: conf.Password,
		DB:       6,
	})
}

package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"log"
)

// 第一个参数和第二个参数分别是接口路径和用户标识，用以唯一确定用户。
// 第三个用户是限流参数，即频次限制
func NewRateLimiter(path string, useSubject bool, limit redis_rate.Limit) func(c *gin.Context) {
	return func(c *gin.Context) {
		key := path
		if useSubject {
			key += ":" + c.GetString("sub")
		}
		ctx := context.Background()
		limiter := redis_rate.NewLimiter(dao.RedisClient)
		res, err := limiter.Allow(ctx, key, limit)
		if err != nil {
			log.Println("[ERROR]: rate_limiter error...")
		}

		log.Println("[Rate Limiter]\tallowed: ", res.Allowed, ", remaining: ", res.Remaining)
		if res.Allowed == 0 {
			panic("[Rate Limiter]\tYou are being rete limited, try again later...")
		}
	}
}

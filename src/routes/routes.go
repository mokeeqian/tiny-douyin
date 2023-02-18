/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/controller"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func InitRouter() *gin.Engine {
	// Create a memory-based store with a rate limit of 1 request per 10 seconds per IP
	rate := limiter.Rate{
		Period: 10 * time.Second,
		Limit:  5,
	}
	store := memory.NewStore()
	limiter := limiter.New(store, rate)

	// Create a Gin middleware using the limiter
	limiterMiddleware := mgin.NewMiddleware(limiter)
	r := gin.Default()
	// 主路由组
	douyinGroup := r.Group("/douyin")
	douyinGroup.Use(limiterMiddleware)
	{
		// user
		userGroup := douyinGroup.Group("/user")
		{
			userGroup.GET("/", middleware.JwtMiddleware(), controller.UserInfo)
			userGroup.POST("/login/", controller.UserLogin)
			userGroup.POST("/register/", controller.UserRegister)
		}

		// publish
		publishGroup := douyinGroup.Group("/publish")
		{
			publishGroup.POST("/action/", middleware.JwtMiddleware(), controller.Publish)
			publishGroup.GET("/list/", middleware.JwtMiddleware(), controller.PublishList)
		}

		// feed
		douyinGroup.GET("/feed/", controller.Feed)

		// favorite
		favoriteGroup := douyinGroup.Group("favorite")
		{
			favoriteGroup.POST("/action/", middleware.JwtMiddleware(), controller.Favorite)
			favoriteGroup.GET("/list/", middleware.JwtMiddleware(), controller.FavoriteList)
		}

		// comment
		commentGroup := douyinGroup.Group("/comment")
		{
			commentGroup.POST("/action/", middleware.JwtMiddleware(), controller.CommentAction)
			commentGroup.GET("/list/", middleware.JwtMiddleware(), controller.CommentList)
		}

		// relation
		relationGroup := douyinGroup.Group("relation")
		{
			relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)
			relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
			relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
			//relationGroup.GET("/friend/list", middleware.JwtMiddleware(), controller.FriendList)
		}

		// message
		messageGroup := douyinGroup.Group("message")
		{
			messageGroup.GET("/chat/", middleware.JwtMiddleware(), controller.MessageChat)
			messageGroup.POST("/action/", middleware.JwtMiddleware(), controller.MessageAction)
		}
	}
	return r
}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis_rate/v9"
	"github.com/mokeeqian/tiny-douyin/src/controller"
	_ "github.com/mokeeqian/tiny-douyin/src/docs" // 这里需要引入本地已生成文档
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	douyinGroup := r.Group("/douyin")
	{
		douyinGroup.GET("/", controller.Welcome)

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
		// 限流，每秒请求一次
		douyinGroup.GET("/feed/", middleware.NewRateLimiter("/douyin/feed/", true, redis_rate.PerSecond(1)), controller.Feed)

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
			relationGroup.GET("/friend/list", middleware.JwtMiddleware(), controller.FriendList)
		}

		// message
		messageGroup := douyinGroup.Group("message")
		{
			messageGroup.GET("/chat/", middleware.JwtMiddleware(), controller.MessageChat)
			messageGroup.POST("/action/", middleware.JwtMiddleware(), controller.MessageAction)
		}
	}

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	return r
}

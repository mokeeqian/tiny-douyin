/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/controller"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	// 主路由组
	douyinGroup := r.Group("/douyin")
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
		//
		//// relation
		//relationGroup := douyinGroup.Group("relation")
		//{
		//	relationGroup.POST("/action/", middleware.JwtMiddleware(), controller.RelationAction)
		//	relationGroup.GET("/follow/list/", middleware.JwtMiddleware(), controller.FollowList)
		//	relationGroup.GET("/follower/list/", middleware.JwtMiddleware(), controller.FollowerList)
		//}
	}

	return r
}

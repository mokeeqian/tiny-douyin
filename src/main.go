/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package main

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/routes"
	"github.com/mokeeqian/tiny-douyin/src/task"
	"github.com/mokeeqian/tiny-douyin/src/util"
)

// @title Tiny-Douyin API
// @version 0.0.1
// @description 短视频社交平台服务端
// @name ssp预备队
// @BasePath /api/v1
func main() {
	//连接数据库
	err := dao.InitMySql()
	if err != nil {
		panic(err)
	}
	//程序退出关闭数据库连接
	defer dao.CloseMysql()

	//连接redis
	dao.InitRedis()

	util.InitFilter()

	// 注册路由
	r := routes.InitRouter()

	// 开启定时任务
	task.CronTaskSetUp()

	dao.SqlSession.AutoMigrate(&db.User{})
	dao.SqlSession.AutoMigrate(&db.Video{})
	dao.SqlSession.AutoMigrate(&db.Comment{})
	dao.SqlSession.AutoMigrate(&db.Favorite{})
	dao.SqlSession.AutoMigrate(&db.Relation{})
	//dao.SqlSession.AutoMigrate(&db.Chat{})
	dao.SqlSession.AutoMigrate(&db.Message{})

	errRun := r.Run(":12138")

	if errRun != nil {
		return
	}

}

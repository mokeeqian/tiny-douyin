/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package main

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/routes"
)

func main() {
	//连接数据库
	err := dao.InitMySql()
	if err != nil {
		panic(err)
	}

	//程序退出关闭数据库连接
	defer dao.Close()

	// 注册路由
	r := routes.InitRouter()

	dao.SqlSession.AutoMigrate(&db.User{})
	dao.SqlSession.AutoMigrate(&db.Video{})
	dao.SqlSession.AutoMigrate(&db.Comment{})
	dao.SqlSession.AutoMigrate(&db.Favorite{})
	dao.SqlSession.AutoMigrate(&db.Follow{})
	dao.SqlSession.AutoMigrate(&db.Message{})

	errRun := r.Run(":12138")

	if errRun != nil {
		return
	}

}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package main

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model"
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

	dao.SqlSession.AutoMigrate(&model.User{})
	dao.SqlSession.AutoMigrate(&model.Video{})
	dao.SqlSession.AutoMigrate(&model.Comment{})
	dao.SqlSession.AutoMigrate(&model.Favorite{})
	dao.SqlSession.AutoMigrate(&model.Follow{})

	errRun := r.Run(":12138")

	if errRun != nil {
		return
	}

}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	DRIVER = "mysql"
)

// SqlSession mysql 连接
var SqlSession *gorm.DB

// InitMySql 初始化连接数据库
func InitMySql() (err error) {
	conf := GetConfig().Mysql

	// 将yaml配置参数拼接成连接数据库的url
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Username,
		conf.Password,
		conf.Url,
		conf.Port,
		conf.Db,
	)

	// 连接数据库
	SqlSession, err = gorm.Open(DRIVER, dsn)
	if err != nil {
		panic(err)
	}

	// 验证数据库连接是否成功，若成功，则无异常
	return SqlSession.DB().Ping()
}

// CloseMysql 关闭数据库连接
func CloseMysql() {
	err := SqlSession.Close()
	if err != nil {
		return
	}
}

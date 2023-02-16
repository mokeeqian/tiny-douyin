/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package dao

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

const (
	DRIVER = "mysql"
)

var SqlSession *gorm.DB

// 连接参数
type conf struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Port     string `yaml:"port"`
}

// 获取配置参数
func (c *conf) getConf() *conf {
	// 配置文件
	yamlFile, err := ioutil.ReadFile("resource/application.yaml")

	if err != nil {
		fmt.Println(err.Error())
	}

	//反序列化
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}

// InitMySql 初始化连接数据库
func InitMySql() (err error) {
	var c conf
	// 获取yaml配置参数
	conf := c.getConf()

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

// Close 关闭数据库连接
func Close() {
	err := SqlSession.Close()
	if err != nil {
		return
	}
}

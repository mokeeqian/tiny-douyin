package dao

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Mysql *MysqlConfig `yaml:"mysql"`
	Redis *RedisConfig `yaml:"redis"`
}

type MysqlConfig struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Port     string `yaml:"port"`
}

type RedisConfig struct {
	Url      string `yaml:"url"`
	Password string `yaml:"password"`
	Db       int    `yaml:"db"`
	Port     string `yaml:"port"`
}

func GetConfig() *Config {
	yamlFile, err := ioutil.ReadFile("resource/application.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	var _config *Config
	err = yaml.Unmarshal(yamlFile, &_config)
	if err != nil {
		fmt.Println(err.Error())
	}
	return _config
}

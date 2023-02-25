package util

import (
	"github.com/importcjj/sensitive"
	"log"
)

var Filter *sensitive.Filter

const sensitiveDict = "./resource/dic/sensitiveDict.txt"

func InitFilter() {
	Filter = sensitive.New()
	err := Filter.LoadWordDict(sensitiveDict)
	if err != nil {
		log.Println("[ERROR]: " + err.Error())
	}
}

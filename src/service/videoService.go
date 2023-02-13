/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"fmt"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"time"
)

// feed 每次返回视频数目
const videoNum = 2

// FeedGet 获得视频列表
func FeedGet(lastTime int64) ([]db.Video, error) {
	//t := time.Now()
	//fmt.Println(t)
	if lastTime == 0 { //没有传入参数或者视屏已经刷完
		lastTime = time.Now().Unix()
	}
	strTime := fmt.Sprint(time.Unix(lastTime, 0).Format("2006-01-02 15:04:05"))
	fmt.Println("查询的截止时间", strTime)
	var VideoList []db.Video
	VideoList = make([]db.Video, 0)
	err := dao.SqlSession.Table("videos").Where("created_at < ?", strTime).Order("created_at desc").Limit(videoNum).Find(&VideoList).Error
	return VideoList, err
}

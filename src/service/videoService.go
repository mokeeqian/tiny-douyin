/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/tencentyun/cos-go-sdk-v5"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// feed 每次返回视频数目
const videoNum = 2

// GetFeed 获得视频列表
func GetFeed(lastTime int64) ([]db.Video, error) {
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

// AddCommentCount add comment_count
func AddCommentCount(videoId uint) error {
	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + 1")).Error; err != nil {
		return err
	}
	return nil
}

// ReduceCommentCount reduce comment_count
func ReduceCommentCount(videoId uint) error {
	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - 1")).Error; err != nil {
		return err
	}
	return nil
}

// GetVideoAuthor get video author
func GetVideoAuthor(videoId uint) (uint, error) {
	var video db.Video
	if err := dao.SqlSession.Table("videos").Where("id = ?", videoId).Find(&video).Error; err != nil {
		return video.ID, err
	}
	return video.AuthorId, nil
}

// AddVideo 添加一条视频信息
func AddVideo(video *db.Video) {
	dao.SqlSession.Table("videos").Create(&video)
}

// GetVideoList 根据用户id查找 所有与该用户相关视频信息
func GetVideoList(userId uint) []db.Video {
	var videoList []db.Video
	dao.SqlSession.Table("videos").Where("author_id=?", userId).Find(&videoList)
	return videoList
}

// CosUpload 上传至云端，返回url
// 该实现参考腾讯云官方GO SDK
func CosUpload(fileName string, reader io.Reader) (string, error) {
	link := "https://qian-1258498110.cos.ap-nanjing.myqcloud.com"
	//u, _ := url.Parse(fmt.Sprintf(dao.COS_URL_FORMAT, dao.COS_BUCKET_NAME, dao.COS_APP_ID, dao.COS_REGION))
	u, _ := url.Parse(link)

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  dao.COS_SECRET_ID,
			SecretKey: dao.COS_SECRET_KEY,
		},
	})
	//path为本地的保存路径
	_, err := client.Object.Put(context.Background(), fileName, reader, nil)
	if err != nil {
		panic(err)
	}
	return link + "/" + fileName, nil
}

// GenerateVideoCover 获取封面
func GenerateVideoCover(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		panic(err)
	}
	return buf
}

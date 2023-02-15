package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	logging "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ReturnAuthor 视频作者结构封装
type ReturnAuthor struct {
	AuthorId      uint   `json:"author_id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// ReturnVideo 创作视频结构封装
type ReturnVideo struct {
	VideoId       uint         `json:"video_id"`
	Author        ReturnAuthor `json:"author"`
	PlayUrl       string       `json:"play_url"`
	CoverUrl      string       `json:"cover_url"`
	FavoriteCount int64        `json:"favorite_count"`
	CommentCount  int64        `json:"comment_count"`
	IsFavorite    bool         `json:"is_favorite"`
	Title         string       `json:"title"`
}

// VideoListResponse 创作列表响应
type VideoListResponse struct {
	common.Response
	VideoList []ReturnVideo `json:"video_list"`
}

// Publish 视频投稿
func Publish(c *gin.Context) {
	//1.校验token
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//2.解析参数
	title := c.PostForm("title")
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}

	//fmt.Println("标题: " + title)

	fileName := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", userId, fileName) // user_id + filename
	//先存储到本地文件夹，再保存到云端，获取封面后最后删除
	saveFile := filepath.Join("../../videos/", finalName)

	// 如果不存在viceos文件夹,创建
	if _, err := os.Stat("../../videos"); os.IsNotExist(err) {
		err := os.Mkdir("../../videos", os.ModePerm)
		if err != nil {
			return
		}
	}

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}

	f, err := data.Open()
	if err != nil {
		err.Error()
	}

	//从本地上传到云端，并获取云端地址
	playUrl, err := service.CosUpload(finalName, f)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  finalName + " 上传成功",
	})

	fmt.Println(playUrl)

	coverName := strings.Replace(finalName, ".mp4", ".jpeg", 1)

	//获取第3帧封面
	img := service.GenerateVideoCover(saveFile, 3)

	//img, _ := jpeg.Decode(buf)//保存到本地时要用到
	//imgw, _ := os.Create(saveImage) //先创建，后写入
	//jpeg.Encode(imgw, img, &jpeg.Options{100})

	// 使用腾讯云,
	//直接传至云端
	coverUrl, err := service.CosUpload(coverName, img)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
		return
	}

	//删除保存在本地中的视频
	err = os.Remove(saveFile) // ignore_security_alert
	if err != nil {
		logging.Info(err)
	}

	//4.保存发布信息至数据库,刚开始发布，喜爱和评论默认为0
	video := db.Video{
		Model:         gorm.Model{},
		AuthorId:      userId,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	}
	// 插入数据库
	service.AddVideo(&video)
}

// PublishList 个人创作列表
func PublishList(c *gin.Context) {
	//1.中间件鉴权token
	getHostId, _ := c.Get("user_id")
	var HostId uint
	if v, ok := getHostId.(uint); ok {
		HostId = v
	}
	//2.查询要查看用户的id的所有视频，返回页面
	getGuestId := c.Query("user_id")
	id, _ := strconv.Atoi(getGuestId)
	GuestId := uint(id)

	// 查询自己的创作列表
	if GuestId == 0 || GuestId == HostId {
		//根据token-id查找用户
		getUser, err := service.GetUser(HostId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				Code: 1,
				Msg:  "Not find this person.",
			})
			c.Abort()
			return
		}

		returnAuthor := ReturnAuthor{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
			IsFollow:      false, // 默认自己不关注自己
		}
		//根据用户id查找 所有相关视频信息
		videoList := service.GetVideoList(HostId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					Code: 1,
					Msg:  "you have no video",
				},
				VideoList: nil,
			})
		} else {
			var returnVideoList []ReturnVideo
			for i := 0; i < len(videoList); i++ {
				curReturnVideo := ReturnVideo{
					VideoId:       videoList[i].ID,
					Author:        returnAuthor,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList = append(returnVideoList, curReturnVideo)
			}
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					Code: 0,
					Msg:  "success",
				},
				VideoList: returnVideoList,
			})
		}
	} else {
		//根据传入id查找用户
		getUser, err := service.GetUser(GuestId)
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				Code: 1,
				Msg:  "Not find this person.",
			})
			c.Abort()
			return
		}

		returnAuthor := ReturnAuthor{
			AuthorId:      getUser.ID,
			Name:          getUser.Name,
			FollowCount:   getUser.FollowCount,
			FollowerCount: getUser.FollowerCount,
			IsFollow:      service.HasRelation(HostId, GuestId),
		}
		//根据用户id查找 所有相关视频信息
		videoList := service.GetVideoList(GuestId)
		if len(videoList) == 0 {
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					Code: 1,
					Msg:  "the one have no video",
				},
				VideoList: nil,
			})
		} else { //需要展示的列表信息
			var returnVideoList []ReturnVideo
			for i := 0; i < len(videoList); i++ {
				curReturnVideo := ReturnVideo{
					VideoId:       videoList[i].ID,
					Author:        returnAuthor,
					PlayUrl:       videoList[i].PlayUrl,
					CoverUrl:      videoList[i].CoverUrl,
					FavoriteCount: videoList[i].FavoriteCount,
					CommentCount:  videoList[i].CommentCount,
					IsFavorite:    service.CheckFavorite(HostId, videoList[i].ID),
					Title:         videoList[i].Title,
				}
				returnVideoList = append(returnVideoList, curReturnVideo)
			}
			c.JSON(http.StatusOK, VideoListResponse{
				Response: common.Response{
					Code: 0,
					Msg:  "success",
				},
				VideoList: returnVideoList,
			})
		}
	}
}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
	"time"
)

/**
专用结构封装
TODO: 这些复杂的结构封装，可以单独抽离出来，放入pack之下
*/

// FeedResponse 视频推荐流响应
type FeedResponse struct {
	common.Response
	VideoList []FeedVideo `json:"video_list,omitempty"`
	NextTime  uint        `json:"next_time,omitempty"`
}

// FeedNoVideoResponse 无视频feed响应
type FeedNoVideoResponse struct {
	common.Response
	NextTime uint `json:"next_time"`
}

// FeedVideo 视频结构封装
type FeedVideo struct {
	Id            uint     `json:"id"`
	Author        FeedUser `json:"author"`
	PlayUrl       string   `json:"play_url"`
	CoverUrl      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}

// FeedUser 作者结构封装
type FeedUser struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// Feed 不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
func Feed(c *gin.Context) {

	strToken := c.Query("token")
	var haveToken bool
	if strToken == "" {
		haveToken = false
	} else {
		haveToken = true
	}
	var strLastTime = c.Query("latest_time")
	lastTime, err := strconv.ParseInt(strLastTime, 10, 32)
	if err != nil {
		lastTime = 0
	}

	var feedVideoList []FeedVideo
	feedVideoList = make([]FeedVideo, 0)
	videoList, _ := service.FeedGet(lastTime)
	var newTime int64 = 0 //返回的视频的最久的一个的时间
	for _, x := range videoList {
		var tmp FeedVideo
		tmp.Id = x.ID
		tmp.PlayUrl = x.PlayUrl
		//tmp.Author = //依靠用户信息接口查询
		var user, err = service.GetUser(x.AuthorId)
		var feedUser FeedUser
		if err == nil { //用户存在
			feedUser.Id = user.ID
			feedUser.FollowerCount = user.FollowerCount
			feedUser.FollowCount = user.FollowCount
			feedUser.Name = user.Name
			//add
			//feedUser.TotalFavorited = user.TotalFavorited
			//feedUser.FavoriteCount = user.FavoriteCount
			feedUser.IsFollow = false
			if haveToken {
				// 查询是否关注
				tokenStruct, ok := middleware.CheckToken(strToken)
				if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
					var uid1 = tokenStruct.UserId //用户id
					var uid2 = x.AuthorId         //视频发布者id
					if service.IsFollowing(uid1, uid2) {
						feedUser.IsFollow = true
					}
				}
			}
		}
		tmp.Author = feedUser
		tmp.CommentCount = x.CommentCount
		tmp.CoverUrl = x.CoverUrl
		tmp.FavoriteCount = x.FavoriteCount
		tmp.IsFavorite = false
		if haveToken {
			//查询是否点赞过
			tokenStruct, ok := middleware.CheckToken(strToken)
			if ok && time.Now().Unix() <= tokenStruct.ExpiresAt { //token合法
				var uid = tokenStruct.UserId         //用户id
				var vid = x.ID                       // 视频id
				if service.CheckFavorite(uid, vid) { //有点赞记录
					tmp.IsFavorite = true
				}
			}
		}
		tmp.Title = x.Title
		feedVideoList = append(feedVideoList, tmp)
		newTime = x.CreatedAt.Unix()
	}
	if len(feedVideoList) > 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "feed获取成功",
			}, //成功
			VideoList: feedVideoList,
			NextTime:  uint(newTime),
		})
	} else {
		c.JSON(http.StatusOK, FeedNoVideoResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "feed未获取到视频",
			}, //成功
			NextTime: 0, //重新循环
		})
	}

}

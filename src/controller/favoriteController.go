package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
)

type FavoriteAuthor struct { //从user中获取,getUser函数
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"` //从following或follower中获取
}

type FavoriteVideo struct { //从video中获取
	Id            uint           `json:"id"`
	Author        FavoriteAuthor `json:"author"`
	PlayUrl       string         `json:"play_url"`
	CoverUrl      string         `json:"cover_url"`
	FavoriteCount int64          `json:"favorite_count"`
	CommentCount  int64          `json:"comment_count"`
	IsFavorite    bool           `json:"is_favorite"` //true
	Title         string         `json:"title"`
}

type FavoriteListResponse struct {
	common.Response
	VideoList []FavoriteVideo `json:"video_list"`
}

func Favorite(c *gin.Context) {
	//user_id获取
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	//参数解析
	actionTypeStr := c.Query("action_type") // 1-点赞，2-取消点赞
	actionType, _ := strconv.ParseUint(actionTypeStr, 10, 10)
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	//函数调用及响应
	err := service.FavoriteAction(userId, uint(videoId), uint(actionType))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			Code: 0,
			Msg:  "点赞成功！",
		})
	}
}

// FavoriteList 获取列表方法
func FavoriteList(c *gin.Context) {
	//user_id获取
	getUserId, _ := c.Get("user_id")
	var userIdHost uint
	if v, ok := getUserId.(uint); ok {
		userIdHost = v
	}
	userIdStr := c.Query("user_id") //自己id或别人id
	userId, _ := strconv.ParseUint(userIdStr, 10, 10)
	userIdNew := uint(userId)
	if userIdNew == 0 {
		userIdNew = userIdHost
	}

	//函数调用及响应
	originalVideoList, err := service.FavoriteList(userIdNew)
	returnVideoList := make([]FavoriteVideo, 0)
	for _, video := range originalVideoList {
		var author = FavoriteAuthor{}
		var getAuthor = db.User{}
		getAuthor, err := service.GetUser(video.AuthorId) //参数类型、错误处理
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
				Code: 403,
				Msg:  "找不到作者！",
			})
			c.Abort()
			return
		}
		//isfollowing
		isfollowing := service.IsFollowing(userIdHost, video.AuthorId) //参数类型、错误处理
		//isfavorite
		isfavorite := service.CheckFavorite(userIdHost, video.ID)
		//作者信息
		author.Id = getAuthor.ID
		author.Name = getAuthor.Name
		author.FollowCount = getAuthor.FollowCount
		author.FollowerCount = getAuthor.FollowerCount
		author.IsFollow = isfollowing
		//组装
		var returnVideo = FavoriteVideo{}
		returnVideo.Id = video.ID //类型转换
		returnVideo.Author = author
		returnVideo.PlayUrl = video.PlayUrl
		returnVideo.CoverUrl = video.CoverUrl
		returnVideo.FavoriteCount = video.FavoriteCount
		returnVideo.CommentCount = video.CommentCount
		returnVideo.IsFavorite = isfavorite
		returnVideo.Title = video.Title

		returnVideoList = append(returnVideoList, returnVideo)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, FavoriteListResponse{
			Response: common.Response{
				Code: 1,
				Msg:  "查找列表失败！",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, FavoriteListResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "已找到列表！",
			},
			VideoList: returnVideoList,
		})
	}
}

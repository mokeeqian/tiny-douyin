package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
)

type RelationUser struct {
	Id            uint   `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

type RelationListResponse struct {
	common.Response
	UserList []RelationUser `json:"user_list"`
}

// RelationAction 登录用户对其他用户进行关注或取消关注
func RelationAction(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	// from id
	fromId := tokenStruct.UserId

	// to id
	toIdInt, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	toId := uint(toIdInt)

	// action type // 1-关注，2-取消关注
	actionTypeInt, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	actionType := uint(actionTypeInt)

	if fromId == toId {
		c.JSON(http.StatusOK, common.Response{
			Code: 405,
			Msg:  "不能关注自己",
		})
		c.Abort()
		return
	}

	// 关注/取关
	err := service.FollowAction(fromId, toId, actionType)

	if err != nil {
		c.JSON(http.StatusBadRequest, common.Response{
			Code: 1,
			Msg:  err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			Code: 0,
			Msg:  "关注/取消关注成功！",
		})
	}
}

// FollowList 关注列表
func FollowList(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)

	// from id
	fromId := tokenStruct.UserId

	// to id
	toIdInt, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	toId := uint(toIdInt)

	var err error
	var tmpFollowList []db.User
	// 查自己
	if toId == 0 {
		tmpFollowList, err = service.FollowList(fromId)
	} else {
		// 查对方
		tmpFollowList, err = service.FollowList(toId)
	}

	// 对返回列表二次加工
	returnFollowList := make([]RelationUser, len(tmpFollowList))
	for i, u := range tmpFollowList {
		returnFollowList[i].Id = u.ID
		returnFollowList[i].Name = u.Name
		returnFollowList[i].FollowCount = u.FollowCount
		returnFollowList[i].FollowerCount = u.FollowerCount
		// 主要计算 关注关系
		returnFollowList[i].IsFollow = service.HasRelation(fromId, u.ID)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, RelationListResponse{
			Response: common.Response{
				Code: 1,
				Msg:  "查询关注列表失败",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "查询关注列表成功",
			},
			UserList: returnFollowList,
		})
	}
}

// FollowerList 粉丝列表
func FollowerList(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)

	// from id
	fromId := tokenStruct.UserId

	// to id
	toIdInt, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	toId := uint(toIdInt)

	var err error
	var tmpFollowerList []db.User
	// 查自己
	if toId == 0 {
		tmpFollowerList, err = service.FollowerList(fromId)
	} else {
		// 查对方
		tmpFollowerList, err = service.FollowerList(toId)
	}

	// 对返回列表二次加工
	returnFollowerList := make([]RelationUser, len(tmpFollowerList))
	for i, u := range tmpFollowerList {
		returnFollowerList[i].Id = u.ID
		returnFollowerList[i].Name = u.Name
		returnFollowerList[i].FollowCount = u.FollowCount
		returnFollowerList[i].FollowerCount = u.FollowerCount
		// 主要计算 关注关系
		returnFollowerList[i].IsFollow = service.HasRelation(fromId, u.ID)
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, RelationListResponse{
			Response: common.Response{
				Code: 1,
				Msg:  "查询粉丝列表失败",
			},
			UserList: nil,
		})
	} else {
		c.JSON(http.StatusOK, RelationListResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "查询粉丝列表成功",
			},
			UserList: returnFollowerList,
		})
	}
}

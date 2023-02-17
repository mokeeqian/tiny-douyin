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

type FriendUser struct {
	RelationUser
	Avatar        string `json:"avatar"`   // 头像 url，暂时写死
	LatestMessage string `json:"message"`  // 和该好友的最新聊天消息
	MessageType   int64  `json:"msg_type"` // message消息的类型，0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
}

type RelationListResponse struct {
	common.Response
	UserList []RelationUser `json:"user_list"`
}

type FriendListResponse struct {
	common.Response
	FriendList []FriendUser `json:"user_list"`
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

// FriendList 好友列表
// 注册登录后，点击消息页面，会立即请求该接口
// 获取可聊天朋友列表，并且会带着和该用户的最新的一条消息
func FriendList(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)

	// from id
	fromId := tokenStruct.UserId

	tmpFriendList, err := service.FriendList(fromId)

	if err != nil {
		c.JSON(http.StatusBadRequest, FriendListResponse{
			Response: common.Response{
				Code: 1,
				Msg:  "查询朋友列表失败",
			},
			FriendList: nil,
		})
	} else {
		// 对返回列表二次加工
		var returnFriendList []FriendUser
		for _, u := range tmpFriendList {
			var msg string
			var msgType int64
			latestMsg1, msgType1, err1 := service.GetLatestMessage(fromId, u.ID)
			latestMsg2, msgType2, err2 := service.GetLatestMessage(u.ID, fromId)
			if err1 != nil && err2 != nil {
				msg = ""
				msgType = -1
			} else if err1 != nil {
				msg = latestMsg2.Content
				msgType = msgType2
			} else if err2 != nil {
				msg = latestMsg1.Content
				msgType = msgType1
			} else {
				if latestMsg1.CreateTime.After(latestMsg2.CreateTime) {
					msg = latestMsg1.Content
					msgType = msgType1
				} else {
					msg = latestMsg2.Content
					msgType = msgType2
				}
			}
			curFriend := FriendUser{
				RelationUser: RelationUser{
					Id:            u.ID,
					Name:          u.Name,
					FollowCount:   u.FollowCount,
					FollowerCount: u.FollowerCount,
					IsFollow:      service.HasRelation(fromId, u.ID),
				},
				Avatar: "https://qian-1258498110.cos.ap-nanjing.myqcloud.com/R.jpg", // 暂时写死
				// TODO: 优化最近一条消息查询
				LatestMessage: msg,
				MessageType:   msgType, // 0 => 当前请求用户接收的消息， 1 => 当前请求用户发送的消息
			}
			returnFriendList = append(returnFriendList, curFriend)
		}
		c.JSON(http.StatusOK, FriendListResponse{
			Response: common.Response{
				Code: 0,
				Msg:  "查询朋友列表成功",
			},
			FriendList: returnFriendList,
		})
	}
}

// CommonFollowList 共同关注
func CommonFollowList(c *gin.Context) {

}

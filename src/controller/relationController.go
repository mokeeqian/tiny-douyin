package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
)

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

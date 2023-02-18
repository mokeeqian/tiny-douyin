package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
)

type MessageResponse struct {
	ID         int    `json:"id"`
	ToUserID   int    `json:"to_user_id"`
	FromUserID int    `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int    `json:"create_time"`
}

// MessageChatResponse 注册响应
type MessageChatResponse struct {
	common.Response
	MessageList []db.Message `json:"message_list"`
}

// MessageChat 详细聊天页面
// 点击朋友列表中的任意用户，进入详细聊天页面。在该页面下会定时轮询消息查询接口，获取最新消息列表
func MessageChat(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	// from id
	fromId := tokenStruct.UserId

	// to id
	toIdInt, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	toId := uint(toIdInt)

	messages, err := service.GetAllMessage(uint(fromId), uint(toId))
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 403,
			Msg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, MessageChatResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "Successfully obtained the messages list.",
		},
		MessageList: messages,
	})
}

// MessageAction 消息发送
func MessageAction(c *gin.Context) {

}

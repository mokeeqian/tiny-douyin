package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Id         uint   `json:"id"`
	ToUserId   uint   `json:"to_user_id"`
	FromUserId uint   `json:"from_user_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"` // 客户端接口是long，db 字段是 datetime
}

type MessageChatResponse struct {
	common.Response
	MessageList []Message `json:"message_list"`
	PreMsgTime  int64     `json:"pre_msg_time"`
}

// MessageChat 详细聊天页面
// 点击朋友列表中的任意用户，进入详细聊天页面。在该页面下客户端会定时轮询消息查询接口，获取最新消息列表
// 服务端需要对消息去重（由pre_msg_time实现）
func MessageChat(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)

	// from id
	fromId := tokenStruct.UserId
	// to id
	toId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	// 上次消息时间
	var preMsgTime int64
	preMsgTimeStr := c.Query("pre_msg_time")
	if preMsgTimeStr == "" {
		preMsgTime = 1546926630
	} else {
		preMsgTime, _ = strconv.ParseInt(preMsgTimeStr, 10, 64)
	}

	msgList, err := service.GetLatestMessageAfter(fromId, uint(toId), preMsgTime)
	// 无消息
	if err != nil {
		c.JSON(http.StatusOK, MessageChatResponse{
			Response: common.Response{
				Code: 1,
				Msg:  "no message",
			},
			MessageList: nil,
			PreMsgTime:  1546926630,
		})
	}
	var responseMsgList []Message
	for _, message := range msgList {
		curMsg := Message{
			Id:         message.Id,
			ToUserId:   message.ToUserId,
			FromUserId: message.FromUserId,
			Content:    message.Content,
			CreateTime: message.CreateTime.Unix(), // 以秒为时间单位
		}
		responseMsgList = append(responseMsgList, curMsg)
	}
	var nextPreMsgTime int64
	if len(responseMsgList) == 0 {
		nextPreMsgTime = 1546926630
	} else {
		nextPreMsgTime = responseMsgList[0].CreateTime
	}
	c.JSON(http.StatusOK, MessageChatResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "get message successfully",
		},
		MessageList: responseMsgList,
		PreMsgTime:  nextPreMsgTime,
	})
}

// MessageAction 消息发送
func MessageAction(c *gin.Context) {
	// 取 token
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)

	// from id
	fromId := tokenStruct.UserId
	// to id
	toId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)

	content := c.Query("content")

	message := db.Message{
		FromUserId: fromId,
		ToUserId:   uint(toId),
		Content:    content,
		CreateTime: time.Now(),
		State:      0,
	}
	err := service.AddMessage(message)
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 1,
			Msg:  "message send failed!",
		})
	} else {
		c.JSON(http.StatusOK, common.Response{
			Code: 0,
			Msg:  "message send successfully!",
		})
	}
}

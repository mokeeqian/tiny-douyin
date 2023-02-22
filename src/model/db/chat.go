package db

import "github.com/jinzhu/gorm"

// Chat 聊天框
type Chat struct {
	gorm.Model
	UserFooId     uint   `json:"user_foo_id"`    // 用户 a
	UserBarId     uint   `json:"user_bar_id"`    // 用户 b
	LatestMessage string `json:"latest_message"` // 最近一条消息的内容
}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import (
	"time"
)

// Message 消息记录
type Message struct {
	Id         uint      `gorm:"primary_key"`
	CreateTime time.Time `json:"create_time"`
	FromUserId uint      `json:"from_user_id"`
	ToUserId   uint      `json:"to_user_id"`
	Content    string    `json:"content"`
	State      int       `json:"state"` // 0: 未读， 1:已读
}

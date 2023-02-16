/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

// Message 消息记录
type Message struct {
	gorm.Model
	FromUserId uint   `json:"from_user_id"`
	ToUserId   uint   `json:"to_user_id"`
	Content    string `json:"content"`
	State      bool   `json:"state"` // 0: 未读， 1:已读
}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

// Relation 好友关注/粉丝关系
type Relation struct {
	gorm.Model
	FromUserId uint `json:"from_user_id"`
	ToUserId   uint `json:"to_user_id"`
	State      uint `json:"state"` // 1 有效， 0 无效
}

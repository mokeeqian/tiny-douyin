/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package model

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username      string `json:"username"`
	Password      string `json:"password"`
	FollowCount   int64  `json:"follow_count"`   // 关注
	FollowerCount int64  `json:"follower_count"` // 粉丝
}

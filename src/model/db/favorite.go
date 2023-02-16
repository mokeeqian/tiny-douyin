/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	VideoId uint `json:"video_id"`
	UserId  uint `json:"user_id"`
	State   uint // 0 无效， 1 有效
}
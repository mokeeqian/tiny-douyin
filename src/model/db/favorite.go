/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	VideoId uint `json:"video_id" gorm:"index"`
	UserId  uint `json:"user_id" gorm:"index"`
	State   uint // 0 无效， 1 有效
}

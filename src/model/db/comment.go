/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	VideoId uint   `json:"video_id" gorm:"index"`
	UserId  uint   `json:"user_id" gorm:"index"`
	Content string `json:"content"`
}

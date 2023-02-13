/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package db

import "github.com/jinzhu/gorm"

type Comment struct {
	gorm.Model
	VideoId uint   `json:"video_id"`
	UserId  uint   `json:"user_id"`
	Content string `json:"content"`
}

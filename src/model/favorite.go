/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package model

import "github.com/jinzhu/gorm"

type Favorite struct {
	gorm.Model
	VideoId uint `json:"video_id"`
	UserId  uint `json:"user_id"`
}

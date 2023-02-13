/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package model

import "github.com/jinzhu/gorm"

type Follow struct {
	gorm.Model
	FromUserId uint `json:"from_user_id"`
	ToUserId   uint `json:"to_user_id"`
}

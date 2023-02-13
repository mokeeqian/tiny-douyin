/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/dao"
)

// CheckFavorite 查询某用户是否点赞某视频
func CheckFavorite(uid uint, vid uint) bool {
	var total int
	if err := dao.SqlSession.Table("favorites").
		Where("user_id = ? AND video_id = ? AND state = 1", uid, vid).Count(&total).
		Error; gorm.IsRecordNotFoundError(err) { //没有该条记录
		return false
	}
	if total == 0 {
		return false
	}
	return true
}

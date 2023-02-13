/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model"
)

// IsFollowing fromId是否关注toId
func IsFollowing(fromId uint, toId uint) bool {
	var relationExist = &model.Follow{}
	//判断关注是否存在
	if err := dao.SqlSession.Model(&model.Follow{}).
		Where("from_user_id=? AND to_user_id=?", fromId, toId).
		First(&relationExist).Error; gorm.IsRecordNotFoundError(err) {
		//关注不存在
		return false
	}
	//关注存在
	return true
}

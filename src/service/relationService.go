/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import "github.com/mokeeqian/tiny-douyin/src/model/db"

/**
关注服务
*/

import (
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/dao"
)

const USER_TABLE_NAME = "users"
const FOLLOW_TABLE_NAME = "follows"

// HasRelation fromId 是否关注 toId; toId 是否有 fromId 这个粉丝
func HasRelation(fromId uint, toId uint) bool {
	var total int
	if err := dao.SqlSession.Table("relations").
		Where("from_user_id = ? AND to_user_id = ? AND state = 1", fromId, toId).Count(&total).
		Error; gorm.IsRecordNotFoundError(err) { //没有该条记录
		return false
	}
	if total == 0 {
		return false
	}
	return true
}

// IncreaseFollowCount 增加 id 的关注数
func IncreaseFollowCount(id uint) error {
	if err := dao.SqlSession.Model(&db.User{}).
		Where("id=?", id).
		Update("follow_count", gorm.Expr("follow_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// DecreaseFollowCount 减少 id 的关注数
func DecreaseFollowCount(id uint) error {
	if err := dao.SqlSession.Model(&db.User{}).
		Where("id=?", id).
		Update("follow_count", gorm.Expr("follow_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// IncreaseFollowerCount 增加 id 的粉丝数
func IncreaseFollowerCount(id uint) error {
	if err := dao.SqlSession.Model(&db.User{}).
		Where("id=?", id).
		Update("follower_count", gorm.Expr("follower_count+?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// DecreaseFollowerCount 减少 id 的粉丝数
func DecreaseFollowerCount(id uint) error {
	if err := dao.SqlSession.Model(&db.User{}).
		Where("id=?", id).
		Update("follower_count", gorm.Expr("follower_count-?", 1)).Error; err != nil {
		return err
	}
	return nil
}

// CreateRelation 创建关注
func CreateRelation(fromId uint, toId uint) error {

	relation := db.Relation{
		FromUserId: fromId,
		ToUserId:   toId,
		State:      1,
	}

	//如果没有记录-则新增，如果有了记录-修改State
	var relationExist = &db.Favorite{}
	result := dao.SqlSession.Table("relations").Where("from_user_id = ? AND to_user_id = ?", fromId, toId).First(&relationExist)

	if result.Error != nil { //不存在
		if err := dao.SqlSession.Table("relations").Create(&relation).Error; err != nil { //创建记录
			return err
		}
	} else { //存在
		// 0 无效，1 有效
		if relationExist.State == 0 {
			dao.SqlSession.Table("relations").Where("from_user_id = ? AND to_user_id = ?", fromId, toId).Update("state", 1)
		}
		return nil
	}
	return nil
}

// DeleteRelation 删除关注
func DeleteRelation(fromId uint, toId uint) error {
	relation := db.Relation{
		FromUserId: fromId,
		ToUserId:   toId,
		State:      0,
	}

	//如果没有记录-则新增，如果有了记录-修改State
	var relationExist = &db.Favorite{}
	result := dao.SqlSession.Table("relations").Where("from_user_id = ? AND to_user_id = ?", fromId, toId).First(&relationExist)

	if result.Error != nil { //不存在
		if err := dao.SqlSession.Table("relations").Create(&relation).Error; err != nil { //创建记录
			return err
		}
	} else { //存在
		// 0 无效，1 有效
		if relationExist.State == 1 {
			dao.SqlSession.Table("relations").Where("from_user_id = ? AND to_user_id = ?", fromId, toId).Update("state", 0)
		}
		return nil
	}
	return nil
}

// FollowAction 关注操作
func FollowAction(fromId uint, toId uint, actionType uint) error {
	//创建关注操作
	if actionType == 1 {
		//判断关注是否存在
		if HasRelation(fromId, toId) {
			//关注存在
			return common.ErrorRelationExit
		} else {
			//关注不存在,创建关注(启用事务Transaction)
			err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
				// from 关注 to，from 成为 to 的粉丝
				err := CreateRelation(fromId, toId)
				if err != nil {
					return err
				}
				//增加from_user_id的关注数
				err = IncreaseFollowCount(fromId)
				if err != nil {
					return err
				}
				//增加to_user_id的粉丝数
				err = IncreaseFollowerCount(toId)
				if err != nil {
					return err
				}
				return nil
			})
			if err1 != nil {
				return err1
			}
		}
	}
	if actionType == 2 {
		//判断关注是否存在
		if HasRelation(fromId, toId) {
			//关注存在,删除关注(启用事务Transaction)
			if err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {
				err := DeleteRelation(fromId, toId)
				if err != nil {
					return err
				}
				//减少from_user_id的关注数
				err = DecreaseFollowCount(fromId)
				if err != nil {
					return err
				}
				//减少to_user_id的粉丝数
				err = DecreaseFollowerCount(toId)
				if err != nil {
					return err
				}
				return nil
			}); err1 != nil {
				return err1
			}

		} else {
			//关注不存在
			return common.ErrorRelationNull
		}
	}
	return nil
}

// FollowingList 获取关注表
func FollowingList(id uint) ([]db.User, error) {
	var userList []db.User

	if err := dao.SqlSession.Model(&db.User{}).
		Joins("left join "+FOLLOW_TABLE_NAME+" on "+USER_TABLE_NAME+".id = "+FOLLOW_TABLE_NAME+".to_user_id").
		Where(FOLLOW_TABLE_NAME+".from_user_id=? AND "+FOLLOW_TABLE_NAME+".state = 1", id).
		Scan(&userList).Error; err != nil {
		return userList, nil
	}
	return userList, nil
}

// FollowerList  获取粉丝表
func FollowerList(Id uint) ([]db.User, error) {
	var userList []db.User

	if err := dao.SqlSession.Model(&db.User{}).
		Joins("left join "+FOLLOW_TABLE_NAME+" on "+USER_TABLE_NAME+".id = "+FOLLOW_TABLE_NAME+".to_user_id").
		Where(FOLLOW_TABLE_NAME+".from_user_id=? AND "+FOLLOW_TABLE_NAME+".state = 1", Id).
		Scan(&userList).Error; err != nil {
		return userList, nil
	}
	return userList, nil
}

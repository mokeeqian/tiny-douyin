package service

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
)

// GetLatestMessage 获取 from 与 to 之间 最近的一条消息内容(from 发给 to)
func GetLatestMessage(fromId uint, toId uint) (db.Message, int64, error) {
	var fromTo []db.Message // from 发给 to
	err := dao.SqlSession.Table("messages").Where("from_user_id = ? AND to_user_id = ?", fromId, toId).Order("create_time desc").Limit(1).Find(&fromTo).Error
	if err != nil || len(fromTo) == 0 {
		return db.Message{}, -1, err
	} else {
		return fromTo[0], 1, nil
	}
}

// GetLatestMessageAfter 获取 from 和 to 之间 在 preMsgTime之后的最近消息
func GetLatestMessageAfter(fromId uint, toId uint, preMsgTime int64) ([]db.Message, error) {
	var msgList []db.Message
	err := dao.SqlSession.Table("messages").Where("(from_user_id = ? AND to_user_id = ? ) OR (from_user_id = ? AND to_user_id = ?)", fromId, toId, toId, fromId).Order("create_time asc").Where("create_time > ?", preMsgTime).Find(&msgList).Error

	return msgList, err
}

func AddMessage(message db.Message) error {
	if err := dao.SqlSession.Table("messages").Create(&message).Error; err != nil {
		return err
	}
	return nil
}

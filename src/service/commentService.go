package service

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"time"
)

// GetCommentListByVideoId 获取指定videoId的评论表
func GetCommentListByVideoId(videoId uint) ([]db.Comment, error) {
	var commentList []db.Comment
	if err := dao.SqlSession.Table("comments").Where("video_id=?", videoId).Find(&commentList).Error; err != nil {
		return commentList, err
	}
	return commentList, nil
}

// PostComment 发布评论
func PostComment(comment db.Comment) error {
	if err := dao.SqlSession.Table("comments").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCommentById 删除指定commentId的评论
func DeleteCommentById(commentId uint) error {
	if err := dao.SqlSession.Table("comments").Where("id = ?", commentId).Update("deleted_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}

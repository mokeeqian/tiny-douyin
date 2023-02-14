package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
)

type UserResponse struct {
	ID            uint   `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

type CommentResponse struct {
	ID         uint         `json:"id,omitempty"`
	Content    string       `json:"content,omitempty"`
	CreateDate string       `json:"create_date,omitempty"`
	User       UserResponse `json:"user,omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment CommentResponse `json:"comment,omitempty"`
}

type CommentListResponse struct {
	common.Response
	CommentList []CommentResponse `json:"comment_list,omitempty"`
}

// CommentAction 评论操作
func CommentAction(c *gin.Context) {
	//1 数据处理
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	actionType := c.Query("action_type") // 1-发布评论，2-删除评论
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	// 2 判断评论操作类型：1代表发布评论，2代表删除评论
	//2.1 非合法操作类型
	if actionType != "1" && actionType != "2" {
		c.JSON(http.StatusOK, common.Response{
			Code: 405,
			Msg:  "Unsupported actionType",
		})
		c.Abort()
		return
	}
	//2.2 合法操作类型
	if actionType == "1" { // 发布评论
		text := c.Query("comment_text")
		PostComment(c, userId, text, uint(videoId))
	} else if actionType == "2" { //删除评论
		commentIdStr := c.Query("comment_id")
		commentId, _ := strconv.ParseInt(commentIdStr, 10, 10)
		DeleteComment(c, uint(videoId), uint(commentId))
	}
}

// PostComment 发布评论
func PostComment(c *gin.Context, userId uint, text string, videoId uint) {
	//1 准备数据模型
	newComment := db.Comment{
		VideoId: videoId,
		UserId:  userId,
		Content: text,
	}

	//2 调用service层发布评论并改变评论数量，获取video作者信息
	err1 := dao.SqlSession.Transaction(func(db *gorm.DB) error {

		// FIXME: 这里其实有问题，应该放入事务之中
		if err := service.PostComment(newComment); err != nil {
			return err
		}
		if err := service.AddCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})
	getUser, err2 := service.GetUser(userId)
	videoAuthor, err3 := service.GetVideoAuthor(videoId)

	//3 响应处理
	if err1 != nil || err2 != nil || err3 != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 403,
			Msg:  "Failed to post comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "post the comment successfully",
		},
		Comment: CommentResponse{
			ID:         newComment.ID,
			Content:    newComment.Content,
			CreateDate: newComment.CreatedAt.Format("01-02"),
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, videoAuthor),
			},
		},
	})
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context, videoId uint, commentId uint) {

	//1 调用service层删除评论并改变评论数量，获取video作者信息
	err := dao.SqlSession.Transaction(func(db *gorm.DB) error {
		if err := service.DeleteCommentById(commentId); err != nil {
			return err
		}
		if err := service.ReduceCommentCount(videoId); err != nil {
			return err
		}
		return nil
	})
	//2 响应处理
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 403,
			Msg:  "Failed to delete comment",
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "delete the comment successfully",
	})
}

// CommentList 获取评论表
func CommentList(c *gin.Context) {
	//1 数据处理
	getUserId, _ := c.Get("user_id")
	var userId uint
	if v, ok := getUserId.(uint); ok {
		userId = v
	}
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.ParseUint(videoIdStr, 10, 10)

	//2.调用service层获取指定videoid的评论表
	commentList, err := service.GetCommentListByVideoId(uint(videoId))

	//2.1 评论表不存在
	if err != nil {
		c.JSON(http.StatusOK, common.Response{
			Code: 403,
			Msg:  "Failed to get commentList",
		})
		c.Abort()
		return
	}

	//2.2 评论表存在
	var responseCommentList []CommentResponse
	for i := 0; i < len(commentList); i++ {
		getUser, err1 := service.GetUser(commentList[i].UserId)

		if err1 != nil {
			c.JSON(http.StatusOK, common.Response{
				Code: 403,
				Msg:  "Failed to get commentList.",
			})
			c.Abort()
			return
		}
		responseComment := CommentResponse{
			ID:         commentList[i].ID,
			Content:    commentList[i].Content,
			CreateDate: commentList[i].CreatedAt.Format("01-02"), // mm-dd
			User: UserResponse{
				ID:            getUser.ID,
				Name:          getUser.Name,
				FollowCount:   getUser.FollowCount,
				FollowerCount: getUser.FollowerCount,
				IsFollow:      service.IsFollowing(userId, commentList[i].ID),
			},
		}
		responseCommentList = append(responseCommentList, responseComment)

	}

	//响应返回
	c.JSON(http.StatusOK, CommentListResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "Successfully obtained the comment list.",
		},
		CommentList: responseCommentList,
	})

}

/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/util"
	logging "github.com/sirupsen/logrus"
	"strconv"
	"time"
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

func InsertFavorite(userId uint, videoId uint) {
	var existsLike db.Favorite
	result := dao.SqlSession.Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).First(&existsLike)
	aLike := db.Favorite{
		UserId:  userId,
		VideoId: videoId,
		State:   1,
	}
	// 点赞记录不存在，则插入
	if result.Error == gorm.ErrRecordNotFound {
		dao.SqlSession.Select("user_id", "video_id", "state", "created_at").Create(&aLike)
	} else {
		//点赞记录存在，则更新
		UpdateFavorite(userId, videoId, 1)
	}
}

func UpdateFavorite(userId uint, videoId uint, state int) {
	dao.SqlSession.Model(db.Favorite{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).Updates(map[string]interface{}{
		"state": state,
	})
}

// FavoriteAction 点赞操作
// 采用redis缓存用户点赞列表、视频点赞数目。Redis消息队列+定时任务更新点赞/取消点赞，定时任务 异步落库点赞数目（一致性要求并不是很高）
// FIXME: 这里数据会有不一致的问题
func FavoriteAction(userId uint, videoId uint, actionType uint) (err error) {
	strUserId := strconv.Itoa(int(userId))
	strVideoId := strconv.Itoa(int(videoId))
	userFavoriteVideoKey := util.KeyUserFavoriteVideo(userId)
	//videoFavoriteByUserKey := util.KeyVideoFavoriteByUser(videoId)
	videoFavoriteCountKey := util.KeyVideoFavoriteCount(videoId)

	//1-点赞
	if actionType == 1 {

		// 查询Redis中是否已缓存过该用户的点赞列表
		// 1、已缓存
		if n, err := dao.RedisClient.Exists(dao.Ctx, userFavoriteVideoKey).Result(); n > 0 {
			if err != nil {
				logging.Errorf("方法FavoriteAction执行失败 %v", err)
				return err
			}
			if _, err1 := dao.RedisClient.SAdd(dao.Ctx, userFavoriteVideoKey, videoId).Result(); err != nil {
				logging.Errorf("方法FavoriteAction执行失败 %v", err)
				return err1
			} else {
				// 将点赞/取消点赞 缓存 在redis中，以"strUserId:videoId的形式存储"，按照 时间顺序，定期更新回数据库
				// TODO: 使用 RocketMQ 重构
				dao.RedisClient.LPush(dao.Ctx, "likeAdd", strUserId+":"+strVideoId)

			}
		} else {
			//2 未缓存
			// 从数据库拉取用户的点赞列表,并缓存到redis中
			videoList, _ := FavoriteList(userId)
			for _, video := range videoList {
				if _, err := dao.RedisClient.SAdd(dao.Ctx, userFavoriteVideoKey, video.ID).Result(); err != nil {
					logging.Errorf("方法：favoriteAction执行失败 %v", err)
					// 防止脏读，直接删除缓存
					dao.RedisClient.Del(dao.Ctx, userFavoriteVideoKey)
					return err
				}
			}

			if _, err := dao.RedisClient.Expire(dao.Ctx, userFavoriteVideoKey, time.Minute*5).Result(); err != nil {
				logging.Errorf("方法favoriteAction：设置过期时间失败%v", err)
				dao.RedisClient.Del(dao.Ctx, userFavoriteVideoKey)
				return err
			}
			// 当前视频的点赞信息放入redis
			if _, err := dao.RedisClient.SAdd(dao.Ctx, userFavoriteVideoKey, videoId).Result(); err != nil {
				logging.Errorf("方法：favoriteAction执行失败 %v", err)
				dao.RedisClient.Del(dao.Ctx, userFavoriteVideoKey)
				return err
			} else {
				dao.RedisClient.LPush(dao.Ctx, "likeAdd", strUserId+":"+strVideoId)
			}
		}

		// 查询当前video的点赞数目是否已缓存
		// 1、已缓存
		if n, err := dao.RedisClient.Exists(dao.Ctx, videoFavoriteCountKey).Result(); n > 0 {
			if err != nil {
				logging.Errorf("方法：favoriteAction: 缓存查询video点赞数目执行失败 %v", err)
				return err
			}
			if _, err := dao.RedisClient.Incr(dao.Ctx, videoFavoriteCountKey).Result(); err != nil {
				logging.Errorf("方法favoriteAction: video点赞数目+1执行失败 %v", err)
				return err
			}
		} else {
			//2、未缓存
			count := GetFavoriteCount(videoId)
			if _, err := dao.RedisClient.Set(dao.Ctx, videoFavoriteCountKey, count, 0).Result(); err != nil {
				logging.Errorf("方法favoriteAction:video点赞数目插入执行失败 %v", err)
				// 防止脏读
				dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
				return err
			}

			//if _, err := dao.RedisClient.Expire(dao.Ctx, videoFavoriteCountKey, time.Minute*5).Result(); err != nil {
			//	logging.Errorf("方法favoriteAction：设置过期时间失败%v", err)
			//	dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
			//	return err
			//}
			if _, err := dao.RedisClient.Incr(dao.Ctx, videoFavoriteCountKey).Result(); err != nil {
				logging.Errorf("方法favoriteAction:video点赞数目+1执行失败 %v", err)
				// 防止脏读
				dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
				return err
			}
		}

	} else { //2-取消点赞

		//存在用户
		if n, err := dao.RedisClient.Exists(dao.Ctx, userFavoriteVideoKey).Result(); n > 0 {
			if err != nil {
				logging.Errorf("方法favoriteAction:缓存查询用户ID执行失败 %v", err)
				return err
			}
			if _, err1 := dao.RedisClient.SRem(dao.Ctx, userFavoriteVideoKey, videoId).Result(); err1 != nil {
				logging.Errorf("方法favoriteAction:缓存取消点赞执行失败 %v", err)
				return err1
			} else {
				dao.RedisClient.LPush(dao.Ctx, "likeDel", strUserId+":"+strVideoId)
			}
		} else { //不存在
			// 从数据库拉取最新的点赞列表,并缓存到数据库中
			videoList, _ := FavoriteList(userId)
			for _, value := range videoList {
				if _, err := dao.RedisClient.SAdd(dao.Ctx, userFavoriteVideoKey, value.ID).Result(); err != nil {
					logging.Errorf("方法：favoriteAction取消点赞执行失败 %v", err)
					// 防止脏读
					dao.RedisClient.Del(dao.Ctx, userFavoriteVideoKey)
					return err
				}
			}
			if _, err := dao.RedisClient.Expire(dao.Ctx, userFavoriteVideoKey, time.Minute*5).Result(); err != nil {
				logging.Errorf("方法favoriteAction：设置过期时间失败%v", err)
				dao.RedisClient.Del(dao.Ctx, userFavoriteVideoKey)
				return err
			}
			// 当前视频取消点赞
			if _, err := dao.RedisClient.SRem(dao.Ctx, userFavoriteVideoKey, videoId).Result(); err != nil {
				logging.Errorf("方法：favoriteAction缓存取消点赞执行失败 %v", err)
				return err
			} else {
				dao.RedisClient.LPush(dao.Ctx, "likeDel", strUserId+":"+strVideoId)
			}
		}

		// 查询当前video的点赞数目是否已缓存
		// 1、已缓存
		if n, err := dao.RedisClient.Exists(dao.Ctx, videoFavoriteCountKey).Result(); n > 0 {
			if err != nil {
				logging.Errorf("方法：favoriteAction: 缓存查询video点赞数目执行失败 %v", err)
				return err
			}
			if _, err := dao.RedisClient.Decr(dao.Ctx, videoFavoriteCountKey).Result(); err != nil {
				logging.Errorf("方法favoriteAction: video点赞数目+1执行失败 %v", err)
				return err
			}
		} else {
			//2、未缓存
			count := GetFavoriteCount(videoId)
			if _, err := dao.RedisClient.Set(dao.Ctx, videoFavoriteCountKey, count, 0).Result(); err != nil {
				logging.Errorf("方法favoriteAction:video点赞数目插入执行失败 %v", err)
				// 防止脏读
				dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
				return err
			}

			//if _, err := dao.RedisClient.Expire(dao.Ctx, videoFavoriteCountKey, time.Minute*5).Result(); err != nil {
			//	logging.Errorf("方法favoriteAction：设置过期时间失败%v", err)
			//	dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
			//	return err
			//}
			if _, err := dao.RedisClient.Decr(dao.Ctx, videoFavoriteCountKey).Result(); err != nil {
				logging.Errorf("方法favoriteAction:video点赞数目+1执行失败 %v", err)
				// 防止脏读
				dao.RedisClient.Del(dao.Ctx, videoFavoriteCountKey)
				return err
			}
		}
	}
	return nil
}

// FavoriteList 获取点赞列表
func FavoriteList(userId uint) ([]db.Video, error) {

	//查询当前id用户的所有点赞视频
	var favoriteList []db.Favorite
	videoList := make([]db.Video, 0)
	if err := dao.SqlSession.Table("favorites").Where("user_id=? AND state=?", userId, 1).Find(&favoriteList).Error; err != nil { //找不到记录
		return videoList, nil
	}
	for _, m := range favoriteList {
		var video = db.Video{}
		if err := dao.SqlSession.Table("videos").Where("id=?", m.VideoId).Find(&video).Error; err != nil {
			return nil, err
		}
		videoList = append(videoList, video)
	}
	return videoList, nil
}

//func GetFavoriteVideoListRedisFirst(userId uint) ([]int, error) {
//	strUserId := strconv.Itoa(int(userId))
//	key := util.KeyUserFavoriteVideo(userId)
//	if n, err := dao.RedisClient.Exists(dao.Ctx, key).Result(); n > 0 { //缓存存在
//		if err != nil {
//			logging.Error("方法GetUserFavoriteVideoList: 缓存获取用户喜爱列表失败%v", err)
//			return nil, err
//		}
//		if strVideoIdList, err := dao.RedisClient.SMembers(dao.Ctx, key).Result(); err != nil {
//			logging.Error("方法GetUserFavoriteVideoList: 缓存获取用户喜爱列表失败%v", err)
//			return nil, err
//		} else {
//			videoIdList := utils.String2Int(strVideoIdList)
//			return videoIdList, nil
//		}
//	} else { //缓存不存在
//		// 从数据库查询，并加载到缓存中
//		videoIdList := like.GetFavoriteVideoIdList(userId)
//		for _, value := range videoIdList {
//			if _, err := dao.RedisClient.SAdd(dao.Ctx, key, value).Result(); err != nil {
//				logging.Error("方法GetUserFavoriteVideoList: 用户喜爱列表加载入缓存失败%v\", err")
//				dao.RedisClient.Del(dao.Ctx, strUserId)
//				return nil, err
//			}
//		}
//		if _, err := dao.RedisClient.Expire(dao.Ctx, strUserId, time.Minute*5).Result(); err != nil {
//			logging.Error("方法favoriteAction：设置过期时间失败%v", err)
//			dao.RedisClient.Del(dao.Ctx, strUserId)
//			return nil, err
//		}
//		return videoIdList, nil
//
//	}
//}

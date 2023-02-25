package util

import "strconv"

const (
	SPLIT           = ":"
	PREFIX_FOLLOWEE = "followee" // 关注
	PREFIX_FOLLOWER = "follower" // 粉丝

	USER_FAVORITE_VIDEO    = "user_favorite_video" // 点赞
	VIDEO_FAVORITE_BY_USER = "video_favorite_by_user"
	VIDEO_FAVORITE_COUNT   = "video_favorite_count"
)

func KeyUserFavoriteVideo(userId uint) string {
	return USER_FAVORITE_VIDEO + SPLIT + strconv.Itoa(int(userId))
}

func KeyVideoFavoriteByUser(videoId uint) string {
	return VIDEO_FAVORITE_BY_USER + SPLIT + strconv.Itoa(int(videoId))
}

func KeyVideoFavoriteCount(videoId uint) string {
	return VIDEO_FAVORITE_COUNT + SPLIT + strconv.Itoa(int(videoId))
}

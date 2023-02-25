package service

import (
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"testing"
)

func TestSaveRedisFavoriteCountToMysql(t *testing.T) {
	dao.InitRedis()
	SaveRedisFavoriteCountToMysql()
}

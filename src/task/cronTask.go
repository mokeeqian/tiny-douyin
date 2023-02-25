package task

import (
	"github.com/henrylee2cn/goutil/calendar/cron"
	"github.com/mokeeqian/tiny-douyin/src/service"
)

func CronTaskSetUp() {
	c := cron.New()
	err := c.AddJob("0/30 * * * * ?", service.SaveRedisFavoriteToMysqlJob{
		Name: "save favorite to db",
	})
	err = c.AddJob("0/30 * * * * ?", service.SaveFavoriteCountToMysqlJob{
		Name: "save favorite count to db",
	})
	if err != nil {
		return
	}
	go c.Start()
	defer c.Stop()
}

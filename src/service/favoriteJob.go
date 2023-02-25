package service

type SaveRedisFavoriteToMysqlJob struct {
	Name string
}

type SaveFavoriteCountToMysqlJob struct {
	Name string
}

// Run 将点赞记录从redis定时刷新到MySQL中
func (saveRedisFavoriteToMysqlJob SaveRedisFavoriteToMysqlJob) Run() {
	SaveRedisFavoriteToMysql()
}

func (saveFavoriteCountToMysqlJob SaveFavoriteCountToMysqlJob) Run() {

}

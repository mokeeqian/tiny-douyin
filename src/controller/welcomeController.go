package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"net/http"
)

func Welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, common.Response{
		Code: 0,
		Msg:  "Welcome to ssp预备队's project...",
	})
}

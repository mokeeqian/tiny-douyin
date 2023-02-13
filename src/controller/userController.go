/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
)

/*
*

	ID Token 响应
*/
type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

/*
*

	注册响应
*/
type UserRegisterResponse struct {
	common.Response
	UserIdTokenResponse
}

// 注册接口
func UserRegister(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	response, err := DoRegister(username, password)

	// 返回响应
	if err != nil {
		c.JSON(http.StatusOK, UserRegisterResponse{
			Response: common.Response{
				Code: 1,
				Msg:  err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, UserRegisterResponse{
		Response: common.Response{
			Code: 0,
		},
		UserIdTokenResponse: response,
	})
	return
}

// 内部逻辑
func DoRegister(username string, password string) (UserIdTokenResponse, error) {
	//0.数据准备
	var userResponse = UserIdTokenResponse{}

	//1.合法性检验
	err := service.IsUserLegal(username, password)
	if err != nil {
		return userResponse, err
	}

	//2.新建用户
	newUser, err := service.CreateUser(username, password)
	if err != nil {
		return userResponse, err
	}

	//3.颁发token
	token, err := middleware.CreateToken(newUser.ID, newUser.Username)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserIdTokenResponse{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

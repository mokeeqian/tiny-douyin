/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/middleware"
	"github.com/mokeeqian/tiny-douyin/src/model/db"
	"github.com/mokeeqian/tiny-douyin/src/service"
	"net/http"
	"strconv"
)

// UserIdTokenResponse ID Token 响应
type UserIdTokenResponse struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

// UserRegisterResponse 注册响应
type UserRegisterResponse struct {
	common.Response
	UserIdTokenResponse
}

// UserLoginResponse 登录响应
type UserLoginResponse struct {
	common.Response
	UserIdTokenResponse
}

// UserInfoQueryResponse
type UserInfoQueryResponse struct {
	UserId        uint   `json:"id"`
	UserName      string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
	//TotalFavorited uint   `json:"total_favorited"`
	//FavoriteCount  uint   `json:"favorite_count"`
}

type UserInfoResponse struct {
	common.Response
	UserList UserInfoQueryResponse `json:"user"`
}

// UserRegister 注册接口
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
			Msg:  "注册成功",
		},
		UserIdTokenResponse: response,
	})
	return
}

// DoRegister 内部逻辑
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
	token, err := middleware.CreateToken(newUser.ID, newUser.Name)
	if err != nil {
		return userResponse, err
	}

	userResponse = UserIdTokenResponse{
		UserId: newUser.ID,
		Token:  token,
	}
	return userResponse, nil
}

// UserLogin 登录
func UserLogin(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	response, err := DoLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{
				Code: 1,
				Msg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "登录成功",
		},
		UserIdTokenResponse: response,
	})
}

// DoLogin 登录逻辑
func DoLogin(username string, password string) (UserIdTokenResponse, error) {
	var response = UserIdTokenResponse{}
	// 用户名密码合法性校验
	err := service.IsUserLegal(username, password)
	if err != nil {
		return response, err
	}

	// 用户存在性校验
	// FIXME: 这里的返回msg有点问题
	var login db.User
	err = service.IsUserExist(username, password, &login)
	if err != nil {
		return response, err
	}

	// 颁发token
	token, err := middleware.CreateToken(login.Model.ID, login.Name)
	if err != nil {
		return response, err
	}

	response = UserIdTokenResponse{
		UserId: login.Model.ID,
		Token:  token,
	}
	return response, nil
}

// UserInfo 用户信息
func UserInfo(c *gin.Context) {
	//根据user_id查询
	rawId := c.Query("user_id")
	userInfoResponse, err := DoInfo(rawId)

	//根据token获得当前用户的userid
	token := c.Query("token")
	tokenStruct, _ := middleware.CheckToken(token)
	hostId := tokenStruct.UserId
	userInfoResponse.IsFollow = service.CheckIsFollow(rawId, hostId)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: common.Response{
				Code: 1,
				Msg:  err.Error(),
			},
		})
		return
	}

	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserInfoResponse{
		Response: common.Response{
			Code: 0,
			Msg:  "查询成功",
		},
		UserList: userInfoResponse,
	})

}

// DoInfo 用户信息
func DoInfo(rawId string) (UserInfoQueryResponse, error) {
	//0.数据准备
	var userInfoQueryResponse = UserInfoQueryResponse{}
	userId, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		return userInfoQueryResponse, err
	}

	//1.获取用户信息
	var user db.User
	err = service.GetUserById(uint(userId), &user)
	if err != nil {
		return userInfoQueryResponse, err
	}

	userInfoQueryResponse = UserInfoQueryResponse{
		UserId:        user.Model.ID,
		UserName:      user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		//TotalFavorited: user.TotalFavorited,
		//FavoriteCount:  user.FavoriteCount,
		IsFollow: false,
	}
	return userInfoQueryResponse, nil
}

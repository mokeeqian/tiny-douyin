/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package service

import (
	"github.com/jinzhu/gorm"
	"github.com/mokeeqian/tiny-douyin/src/common"
	"github.com/mokeeqian/tiny-douyin/src/dao"
	"github.com/mokeeqian/tiny-douyin/src/model"
	"golang.org/x/crypto/bcrypt"
)

const (
	UsernameMaxLength = 32
	PasswordMaxLength = 32
	PasswordMinLength = 8
)

// 功能函数

// HashAndSalt 加密密码
func HashAndSalt(pwdStr string) (pwdHash string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return
	}
	pwdHash = string(hash)
	return
}

// CheckPasswords 验证密码
func CheckPasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}

// IsUserLegal 用户名和密码合法性检验
func IsUserLegal(userName string, passWord string) error {
	//1.用户名检验
	if userName == "" {
		return common.ErrorUserNameNull
	}
	if len(userName) > UsernameMaxLength {
		return common.ErrorUserNameLength
	}
	//2.密码检验
	if passWord == "" {
		return common.ErrorPasswordNull
	}
	if len(passWord) > PasswordMaxLength || len(passWord) < PasswordMinLength {
		return common.ErrorPasswordLength
	}
	return nil
}

// 业务函数

func CreateUser(username string, password string) (model.User, error) {
	// 密码加密
	encryptedPassword, _ := HashAndSalt(password)
	// 创建数据模型
	newUser := model.User{
		Username: username,
		Password: encryptedPassword,
	}
	//2.模型关联到数据库表users //可注释
	dao.SqlSession.AutoMigrate(&model.User{})
	//3.新建user
	if IsUserExistByName(username) {
		//用户已存在
		return newUser, common.ErrorUserExist
	} else {
		//用户不存在，新建用户
		if err := dao.SqlSession.Model(&model.User{}).Create(&newUser).Error; err != nil {
			//错误处理
			//fmt.Println(err)
			panic(err)
			return newUser, err
		}
	}
	return newUser, nil
}

func IsUserExistByName(username string) bool {
	var userExist = &model.User{}
	if err := dao.SqlSession.Model(&model.User{}).Where("username=?", username).First(&userExist).Error; gorm.IsRecordNotFoundError(err) {
		//不存在
		return false
	}
	//存在
	return true
}

func IsUserExist(username string, password string, login *model.User) error {
	if login == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("username=?", username).First(login)
	if !CheckPasswords(login.Password, password) {
		return common.ErrorPasswordFalse
	}
	if login.Model.ID == 0 {
		return common.ErrorAll
	}
	return nil
}

// GetUser 根据用户id获取用户信息
func GetUser(userId uint) (model.User, error) {
	//1.数据模型准备
	var user model.User
	//2.在users表中查对应user_id的user
	if err := dao.SqlSession.Model(&model.User{}).Where("id = ?", userId).Find(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetUserById 根据用户id获取用户信息，用于userInfo
func GetUserById(userId uint, user *model.User) error {
	if user == nil {
		return common.ErrorNullPointer
	}
	dao.SqlSession.Where("id=?", userId).First(user)
	return nil
}

//// CheckIsFollow 检验已登录用户是否关注目标用户
//func CheckIsFollow(targetId string, userid uint) bool {
//	//1.修改targetId数据类型
//	hostId, err := strconv.ParseUint(targetId, 10, 64)
//	if err != nil {
//		return false
//	}
//	//如果是自己查自己，那就是没有关注
//	if uint(hostId) == userid {
//		return false
//	}
//	//2.自己是否关注目标userId
//	return IsFollowing(uint(hostId), userid)
//}

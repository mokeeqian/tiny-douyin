/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package common

import "errors"

/*
*

	错误状态枚举
*/
var (
	ErrorUserNameNull   = errors.New("用户名为空")
	ErrorUserNameLength = errors.New("用户名长度不符合规范")
	ErrorPasswordNull   = errors.New("密码为空")
	ErrorPasswordLength = errors.New("密码长度不符合规范")
	ErrorUserExist      = errors.New("用户已存在")
	ErrorAll            = errors.New("用户不存在，账号或密码出错")
	ErrorNullPointer    = errors.New("空指针异常")
	ErrorPasswordFalse  = errors.New("密码错误")
	ErrorRelationExit   = errors.New("关注已存在")
	ErrorRelationNull   = errors.New("关注不存在")
)

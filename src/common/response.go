/*
 * Copyright (c) 2023.
 * mokeeqian
 */

package common

type Response struct {
	// 这里的json部分保持与 apifox 文档完全一致
	Code int32  `json:"status_code"`
	Msg  string `json:"status_msg,omitempty"`
}

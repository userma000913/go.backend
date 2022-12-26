package model

type (
	// ApiResult 接口响应结果
	ApiResult struct {
		// 错误码
		Code int `json:"code"`
		// 错误提示语
		Msg string `json:"msg"`
	}

	// RespBase 响应基类
	RespBase struct {
		ApiResult
		Data interface{} `json:"data"`
	}
)

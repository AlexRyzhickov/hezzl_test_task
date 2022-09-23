package utils

import "hezzl_test_task/internal/service"

type notFoundMsgStruct struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func NotFoundMsg() interface{} {
	return notFoundMsgStruct{
		Code:    3,
		Message: service.NotFoundError,
		Details: struct{}{},
	}
}

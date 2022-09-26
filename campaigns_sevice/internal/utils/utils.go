package utils

import (
	"github.com/nats-io/nats.go"
	"hezzl_test_task/internal/service"
	"log"
)

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

type Logger struct {
	nc *nats.Conn
}

func InitLogger(nc *nats.Conn) Logger {
	return Logger{nc: nc}
}

func (l Logger) Write(p []byte) (int, error) {
	err := l.nc.Publish("foo", p)
	if err != nil {
		return 0, err
	}
	log.Print(string(p))
	return len(p), nil
}

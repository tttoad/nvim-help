package action

import (
	"encoding/json"
	"unsafe"
)

type Code int64

const (
	SuccessCode Code = 200
	FailCode    Code = 500
)

type Result struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewSuccessResult(data interface{}) *Result {
	if data == nil {
		return &Result{Code: SuccessCode}
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return NewFailResult(err)
	}

	return &Result{
		Code: SuccessCode,
		Data: unsafe.String(&dataBytes[0], len(dataBytes)),
	}
}

func NewFailResult(err error) *Result {
	return &Result{
		Code:    FailCode,
		Message: err.Error(),
	}
}

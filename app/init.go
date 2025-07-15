package app

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/mylukin/EchoPilot-Template/app/middleware"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type (
	JSON   = middleware.JSON
	Error  = middleware.Error
	Result = middleware.Result
)

// Throw is throw err
func Throw(i interface{}, code ...int) *Error {
	return middleware.Throw(i, code...)
}

// Success
func Success(data any) *Result {
	return middleware.Success(data)
}

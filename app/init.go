package app

import (
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {

}

type JSON map[string]interface{}

// Error
func (body JSON) Error() string {
	data, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// Error is error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"-"`
}

func (e *Error) Error() string {
	return e.Message
}

// Result is send message
type Result struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Result any    `json:"result"`
}

// Error
func (body Result) Error() string {
	data, err := json.Marshal(body)
	if err != nil {
		return err.Error()
	}
	return string(data)
}

// Throw is throw err
func Throw(i interface{}, code ...int) *Error {
	errCode := 400
	if len(code) > 0 {
		errCode = code[0]
	}

	var message string
	var data interface{}
	switch v := i.(type) {
	case string:
		message = v
	case Error:
		message = v.Message
		data = v.Data
		if len(code) == 0 && v.Code != 0 {
			errCode = v.Code
		}
	default:
		message = fmt.Sprintf("%v", i)
	}

	return &Error{
		Code:    errCode,
		Message: message,
		Data:    data,
	}
}

// Success
func Success(data any) *Result {
	return &Result{
		Code:   200,
		Status: "ok",
		Result: data,
	}
}

package middleware

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot-Template/app/model"
	"github.com/mylukin/EchoPilot/helper"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type (
	JSON   = model.JSON
	Error  = model.Error
	Result = model.Result
)

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

// ResponseToJSON is response to json
func ResponseToJSON() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			result := next(c)
			c.Logger().Debugf("result type: %T", result)
			// error
			if fmt.Sprintf("%T", result) == "*errors.errorString" {
				err := helper.HiddenBotToken(result.Error())
				// 防止bot token 泄漏
				result = &Error{Message: err}
			}
			// 检查类型
			switch v := result.(type) {
			case *Error:
				code := v.Code
				if code == 0 {
					code = 400
				}
				return c.JSON(http.StatusOK, Result{
					Code:    code,
					Status:  "err",
					Message: helper.HiddenBotToken(v.Error()),
				})
			case *echo.HTTPError:
				code := v.Code
				if code == 0 {
					code = 400
				}
				return c.JSON(http.StatusOK, Result{
					Code:    code,
					Status:  "err",
					Message: helper.HiddenBotToken(fmt.Sprintf("%v", v.Message)),
				})
			default:
				if c.Response().Size == 0 {
					return c.JSON(http.StatusOK, v)
				}
				return result
			}
		}
	}
}

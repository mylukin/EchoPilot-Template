package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot-Template/app"
	"github.com/mylukin/EchoPilot/helper"
)

type (
	// Result is api result
	Result struct {
		Code    int         `json:"code"`
		Status  string      `json:"status"`
		Message string      `json:"message,omitempty"`
		Result  interface{} `json:"result,omitempty"`
	}
)

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
				result = &app.Error{Message: err}
			}
			// 检查类型
			switch v := result.(type) {
			case *app.Error:
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

package middleware

import (
	"net/http"
	"regexp"

	"github.com/golang-jwt/jwt/v5"
	echoJWT "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot-Template/config"
)

// validate token
func ValidateToken() echo.MiddlewareFunc {
	return echoJWT.WithConfig(echoJWT.Config{
		Skipper: func(c echo.Context) bool {
			// 定义要跳过的路径的正则表达式列表
			skipPaths := []string{
				`^/api/v1/config$`, // 精确匹配 /api/v1/config
				// 添加更多路径的正则表达式，例如:
				// `^/api/v1/some-path$`, // 精确匹配 /api/v1/some-path
				// `^/api/v1/public/.*$`, // 匹配以 /api/v1/public/ 开头的任何路径
			}
			// 遍历所有的正则表达式
			for _, pathRegex := range skipPaths {
				matched, _ := regexp.MatchString(pathRegex, c.Path())
				if matched {
					return true // 如果当前路径匹配，则跳过认证
				}
			}

			return false // 如果没有找到匹配，不跳过认证
		},
		ContextKey:  "User",
		TokenLookup: "header:token,cookie:token,query:token",
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(config.UserLoginClaims)
		},
		SigningKey: []byte(config.JWTToken),
		ErrorHandler: func(c echo.Context, err error) error {
			// 定义要跳过的路径的正则表达式列表
			skipPaths := []string{
				// `^/api/v1/test/[^/]+$`, // 匹配 /api/v1/test/ 后面跟任何字母或数字
			}

			// 遍历所有的正则表达式
			for _, pathRegex := range skipPaths {
				matched, _ := regexp.MatchString(pathRegex, c.Path())
				if matched {
					return nil
				}
			}

			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
		ContinueOnIgnoredError: true,
	})
}

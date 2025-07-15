package routers

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot-Template/app"
	appMiddleware "github.com/mylukin/EchoPilot-Template/app/middleware"
	"github.com/mylukin/EchoPilot-Template/config"
	"github.com/mylukin/EchoPilot/middleware"
)

// MountAPI is api router
func MountAPI(e *echo.Echo) {
	api := e.Group("/api", appMiddleware.ResponseToJSON())

	// 设置语言
	api.Use(middleware.SetLang(middleware.SetLangConfig{
		Languages: config.Languages,
		Language: func(req *http.Request, res *echo.Response, c echo.Context) string {
			return req.Header.Get("Accept-Language")
		},
	}))

	// 频率限制
	api.Use(middleware.RateLimiting(&middleware.RateLimitingConfig{
		// 频率限制，单位：毫秒
		Window: 1000,
		// allow request number
		Limit: 10,
		// generate id
		Generator: func(req *http.Request, res *echo.Response, c echo.Context) string {
			requestDump, _ := httputil.DumpRequest(req, true)
			return fmt.Sprintf("%x", md5.Sum(requestDump))
		},
		// 回调
		Callback: func(req *http.Request, res *echo.Response, c echo.Context) error {
			// 返回 429
			return echo.ErrTooManyRequests
		},
	}))

	api.GET("/", app.HelloWorld)
	api.GET("/ping", app.Ping)
	
	// debug pprof
	PProf(api)

}

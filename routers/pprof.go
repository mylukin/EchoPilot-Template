package routers

import (
	"net/http/pprof"

	"github.com/labstack/echo/v4"
)

// PProf is pprof
func PProf(g *echo.Group) {
	g.GET("/debug/", func(c echo.Context) error {
		pprof.Index(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/allocs", func(c echo.Context) error {
		pprof.Handler("allocs").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/heap", func(c echo.Context) error {
		pprof.Handler("heap").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/goroutine", func(c echo.Context) error {
		pprof.Handler("goroutine").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/block", func(c echo.Context) error {
		pprof.Handler("block").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/threadcreate", func(c echo.Context) error {
		pprof.Handler("threadcreate").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/cmdline", func(c echo.Context) error {
		pprof.Cmdline(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/profile", func(c echo.Context) error {
		pprof.Profile(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/symbol", func(c echo.Context) error {
		pprof.Symbol(c.Response().Writer, c.Request())
		return nil
	})
	g.POST("/debug/symbol", func(c echo.Context) error {
		pprof.Symbol(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/trace", func(c echo.Context) error {
		pprof.Trace(c.Response().Writer, c.Request())
		return nil
	})
	g.GET("/debug/mutex", func(c echo.Context) error {
		pprof.Handler("mutex").ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})
}

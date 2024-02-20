package main

//go:generate codetool gen_bot_events github.com/mylukin/EchoPilot-Template
//go:generate easyi18n extract . ./locales/en.json
//go:generate easyi18n update -f ./locales/en.json ./locales/zh-hans.json
//go:generate easyi18n update -f ./locales/en.json ./locales/zh-hant.json
//go:generate easyi18n generate --pkg=catalog ./locales ./catalog/main.go

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mylukin/EchoPilot-Template/config"
	"github.com/mylukin/EchoPilot-Template/routers"
	"github.com/mylukin/EchoPilot/helper"
	eMiddleware "github.com/mylukin/EchoPilot/middleware"
	redisDb "github.com/mylukin/EchoPilot/storage/redis"

	_ "github.com/mylukin/EchoPilot-Template/catalog"
)

const APP_NAME = "EchoPilot"
const APP_VERSION = "0.1.0"

func init() {
	// load .env
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	redisDb.Prefix(config.CachePrefix)
}

func main() {

	e := echo.New()
	// hidden Banner
	e.HideBanner = true
	// debug mode
	e.Debug = helper.Config("ENV") != "GA"
	// enable logger
	e.Use(eMiddleware.LoggerWithConfig(eMiddleware.LoggerConfig{
		Format:      middleware.DefaultLoggerConfig.Format,
		Timeout:     200 * time.Millisecond,
		MinBodySize: 5,
	}))
	// Recover middleware recovers from panics anywhere in the chain, prints stack trace and handles the control to the centralized HTTPErrorHandler.
	e.Use(middleware.Recover())
	// add request id
	e.Use(middleware.RequestID())
	// custom Powered-By
	e.Use(eMiddleware.PoweredBy(eMiddleware.PoweredByConfig{
		Name:    APP_NAME,
		Version: APP_VERSION,
	}))
	// add CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	// Body Dump
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if c.Echo().Debug {
			requestDump, _ := httputil.DumpRequest(c.Request(), true)
			fmt.Printf("request: %s\n\n", requestDump)

			reqContentType := http.DetectContentType(reqBody)
			if strings.Contains(reqContentType, "text/") {
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, reqBody)
			} else {
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, reqContentType, len(reqBody)))
			}
			resContentType := http.DetectContentType(resBody)
			if strings.Contains(reqContentType, "text/") {
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, resBody)
			} else {
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, resContentType, len(resBody)))
			}
		}
	}))

	// static
	e.Static("/static", "public")

	// mount routers
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// mount api
	routers.MountAPI(e)

	// Start server
	go func() {
		if err := e.Start(":" + helper.Config("PORT")); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatalf("shutting down the server: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

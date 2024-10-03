package main

//go:generate codetool gen_bot_events github.com/mylukin/EchoPilot-Template
//go:generate easyi18n extract . ./locales/en.json
//go:generate easyi18n update -f ./locales/en.json ./locales/zh-hans.json
//go:generate translator -m gpt-4o -i ./locales/en.json -l zh-hans
//go:generate translator -m gpt-4o -i ./locales/zh-hans.json -l zh-hant
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

	"github.com/Xuanwo/go-locale"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/mylukin/EchoPilot-Template/app"
	"github.com/mylukin/EchoPilot-Template/command"
	"github.com/mylukin/EchoPilot-Template/config"
	"github.com/mylukin/EchoPilot-Template/routers"
	"github.com/mylukin/EchoPilot/helper"
	eMiddleware "github.com/mylukin/EchoPilot/middleware"
	"github.com/mylukin/EchoPilot/storage/mongo"
	redisDb "github.com/mylukin/EchoPilot/storage/redis"
	ei18n "github.com/mylukin/easy-i18n/i18n"
	"github.com/urfave/cli/v2"

	_ "github.com/mylukin/EchoPilot-Template/catalog"
)

const APP_NAME = "{APP_NAME}"
const APP_VERSION = "0.1.0"

func init() {
	redisDb.Prefix(config.CachePrefix)
}

// see: https://github.com/go-playground/validator
type gValidator struct {
	validator *validator.Validate
}

func (cv *gValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	// 关闭数据库连接
	defer mongo.Close()

	if len(os.Args) > 1 && (os.Args[1] == "main.go" || os.Args[1] == "server") {
		handleHttp()
	} else {
		handleCLI()
	}
}

// CLI模式的处理逻辑
func handleCLI() {
	// Detect OS language
	tag, _ := locale.Detect()

	// Set log level
	if helper.Config("ENV") != "GA" {
		log.SetLevel(log.DEBUG)
	} else {
		log.SetLevel(log.INFO)
	}

	// Set Language
	ei18n.SetLang(tag)

	app := &cli.App{
		Name:  `{APP_NAME}`,
		Usage: ei18n.Sprintf("A new cli application"),
		Action: func(c *cli.Context) error {
			cli.ShowAppHelp(c)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "server",
				Aliases: []string{"s"},
				Usage:   ei18n.Sprintf("Start http server"),
				Action: func(c *cli.Context) error {
					handleHttp()
					return nil
				},
			},
		},
	}

	// 注册命令
	command.RegisterCommands(app)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// HTTP请求的处理逻辑
func handleHttp() {
	e := echo.New()
	// hidden Banner
	e.HideBanner = true
	// debug mode
	e.Debug = helper.Config("ENV") != "GA"
	// Set log level
	if e.Debug {
		log.SetLevel(log.DEBUG)
	} else {
		log.SetLevel(log.INFO)
	}
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

	validateObj := validator.New()
	validateObj.RegisterValidation("media", app.IsURLOrDataURI)
	// set Validator
	e.Validator = &gValidator{validator: validateObj}

	// Body Dump
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		if c.Echo().Debug {
			requestDump, _ := httputil.DumpRequest(c.Request(), true)
			fmt.Printf("request: %s\n\n", requestDump)

			reqContentType := http.DetectContentType(reqBody)
			if strings.Contains(reqContentType, "text/") {
				if len(reqBody) > 1024 {
					reqBody = []byte(fmt.Sprintf("%s...", reqBody[:1024]))
				}
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, reqBody)
			} else {
				fmt.Printf("---- %s %s reqBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, reqContentType, len(reqBody)))
			}
			resContentType := http.DetectContentType(resBody)
			if strings.Contains(resContentType, "text/") {
				if len(resBody) > 1024 {
					resBody = []byte(fmt.Sprintf("%s...", resBody[:1024]))
				}
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, resBody)
			} else {
				fmt.Printf("---- %s %s resBody: %s\n", c.Request().Method, c.Request().RequestURI, fmt.Sprintf(`%v, %v`, resContentType, len(resBody)))
			}
		}
	}))

	// static
	e.Static("/", "public")

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

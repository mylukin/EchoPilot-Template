package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot/service/i18n"
)

// HelloWorld hello for API
func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, i18n.Sprintf(c, "Hello, %s!", "API"))
}

// Ping is ping
func Ping(c echo.Context) error {
	return c.String(http.StatusOK, `ok`)
}

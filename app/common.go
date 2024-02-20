package app

import (
	"bytes"
	"io"

	"github.com/labstack/echo/v4"
	"github.com/mylukin/EchoPilot/helper"
	"github.com/mylukin/easy-i18n/i18n"
)

// get user language
func GetUserLang(c echo.Context) string {
	if c.Get("Language") == nil {
		return helper.Config("LANGUAGE")
	}
	return c.Get("Language").(*i18n.Printer).String()
}

// GetReqBody get request body
func GetReqBody(c echo.Context) []byte {
	reqBody := []byte{}
	if c.Request().Body != nil { // Read
		reqBody, _ = io.ReadAll(c.Request().Body)
	}
	c.Request().Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Reset
	return reqBody
}

// GetCookie get cookie
func GetCookie(c echo.Context, name string) string {
	cookie, err := c.Cookie(name)
	if err != nil {
		return ""
	}
	return cookie.Value
}

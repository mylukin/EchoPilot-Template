package middleware

import "gopkg.in/telebot.v3"

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(telebot.Context) bool

	// BeforeFunc defines a function which is executed just before the middleware.
	BeforeFunc func(telebot.Context)
)

// DefaultSkipper returns false which processes the middleware.
func DefaultSkipper(telebot.Context) bool {
	return false
}

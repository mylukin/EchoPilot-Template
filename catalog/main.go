package catalog

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// init
func init() {
	initEn(language.Make("en"))
	initZhhans(language.Make("zh-hans"))
	initZhhant(language.Make("zh-hant"))
}
// initEn will init en support.
func initEn(tag language.Tag) {
	message.SetString(tag, "A new cli application", "A new cli application")
	message.SetString(tag, "EchoPilot test", "EchoPilot test")
	message.SetString(tag, "Executing test command...", "Executing test command...")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "Start http server", "Start http server")
}
// initZhhans will init zh-hans support.
func initZhhans(tag language.Tag) {
	message.SetString(tag, "A new cli application", "A new cli application")
	message.SetString(tag, "EchoPilot test", "EchoPilot test")
	message.SetString(tag, "Executing test command...", "Executing test command...")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "Start http server", "Start http server")
}
// initZhhant will init zh-hant support.
func initZhhant(tag language.Tag) {
	message.SetString(tag, "A new cli application", "A new cli application")
	message.SetString(tag, "EchoPilot test", "EchoPilot test")
	message.SetString(tag, "Executing test command...", "Executing test command...")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "Start http server", "Start http server")
}

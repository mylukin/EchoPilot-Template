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
	message.SetString(tag, "%s gen_bot_events [module] [outfile]", "%s gen_bot_events [module] [outfile]")
	message.SetString(tag, "Generate Bot Events", "Generate Bot Events")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "[module] can't be empty.", "[module] can't be empty.")
	message.SetString(tag, "a tool for managing message translations.", "a tool for managing message translations.")
}

// initZhhans will init zh-hans support.
func initZhhans(tag language.Tag) {
	message.SetString(tag, "%s gen_bot_events [module] [outfile]", "%s gen_bot_events [module] [outfile]")
	message.SetString(tag, "Generate Bot Events", "Generate Bot Events")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "[module] can't be empty.", "[module] can't be empty.")
	message.SetString(tag, "a tool for managing message translations.", "a tool for managing message translations.")
}

// initZhhant will init zh-hant support.
func initZhhant(tag language.Tag) {
	message.SetString(tag, "%s gen_bot_events [module] [outfile]", "%s gen_bot_events [module] [outfile]")
	message.SetString(tag, "Generate Bot Events", "Generate Bot Events")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "[module] can't be empty.", "[module] can't be empty.")
	message.SetString(tag, "a tool for managing message translations.", "a tool for managing message translations.")
}

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
	message.SetString(tag, "Executing test command...", "Executing test command...")
	message.SetString(tag, "Hello, %s!", "Hello, %s!")
	message.SetString(tag, "Start http server", "Start http server")
	message.SetString(tag, "test", "test")
}
// initZhhans will init zh-hans support.
func initZhhans(tag language.Tag) {
	message.SetString(tag, "A new cli application", "一个新的cli应用程序")
	message.SetString(tag, "Executing test command...", "执行测试命令...")
	message.SetString(tag, "Hello, %s!", "你好, %s!")
	message.SetString(tag, "Start http server", "启动http服务器")
	message.SetString(tag, "test", "测试")
}
// initZhhant will init zh-hant support.
func initZhhant(tag language.Tag) {
	message.SetString(tag, "A new cli application", "一個新的cli應用程式")
	message.SetString(tag, "Executing test command...", "執行測試命令...")
	message.SetString(tag, "Hello, %s!", "你好, %s!")
	message.SetString(tag, "Start http server", "啟動http伺服器")
	message.SetString(tag, "test", "測試")
}

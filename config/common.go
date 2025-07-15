package config

import (
	"github.com/mylukin/EchoPilot/helper"
)

// ENV is environment
var ENV string = helper.Config("ENV")

// Default language
var Language string = helper.Config("LANGUAGE")

// Default timezone
var TZ string = helper.Config("TZ")

// URL
var URL string = helper.Config("URL")

// Fontend URL
var FrontendURL string = helper.Config("FRONTEND_URL")

// BotToken is bot token
var BotToken string = helper.Config("TG_BOT_TOKEN")

// BotUsername is bot username
var BotUsername string = helper.Config("TG_BOT_USERNAME")

// WebAppUsername is mini app username
var WebAppUsername string = BotUsername

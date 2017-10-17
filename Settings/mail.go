package Settings


import (
	util "github.com/horae-app/api/Util"
)

var MAIL_FROM = util.GetEnv("HORAE_APP_MAIL_FROM", "")
var MAIL_USER = util.GetEnv("HORAE_APP_MAIL_USER", "")
var MAIL_PWD = util.GetEnv("HORAE_APP_MAIL_PWD", "")
var MAIL_PORT = util.GetEnv("HORAE_APP_MAIL_PORT", "")
var MAIL_SMTP = util.GetEnv("HORAE_APP_MAIL_SMTP", "")

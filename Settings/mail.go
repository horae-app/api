package Settings

import (
	"strconv"
)

var MAIL_FROM = GetEnv("HORAE_APP_MAIL_FROM", "")
var MAIL_USER = GetEnv("HORAE_APP_MAIL_USER", "")
var MAIL_PWD = GetEnv("HORAE_APP_MAIL_PWD", "")
var MAIL_PORT, _ = strconv.Atoi(GetEnv("HORAE_APP_MAIL_PORT", ""))
var MAIL_SMTP = GetEnv("HORAE_APP_MAIL_SMTP", "")

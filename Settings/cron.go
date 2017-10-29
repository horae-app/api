package Settings

import (
	"strconv"
)

var CONTACT_REMEMBER_HOUR, _ = strconv.Atoi(GetEnv("HORAE_APP_CONTACT_REMEMBER_HOUR", ""))
var CONTACT_REMEMBER_MINUTE, _ = strconv.Atoi(GetEnv("HORAE_APP_CONTACT_REMEMBER_MINUTE", ""))

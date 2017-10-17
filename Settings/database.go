package Settings


import (
	util "github.com/horae-app/api/Util"
)

var DB_ENDPOINT = util.GetEnv("HORAE_APP_DB_ADDRRESS", "127.0.0.1")
var DB_NAME = util.GetEnv("HORAE_APP_DB_NAME", "horaeapi")

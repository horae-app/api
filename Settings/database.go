package Settings

var DB_ENDPOINT = GetEnv("HORAE_APP_DB_ADDRRESS", "127.0.0.1")
var DB_NAME = GetEnv("HORAE_APP_DB_NAME", "horaeapi")

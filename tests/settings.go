package tests

import (
	"fmt"
	"go-final/pkg/config"
)

var Port = config.GetPortFromEnv()
var DBFile = fmt.Sprintf("../%s", config.GetDBFromEnv())
var FullNextDate = true
var Search = false
var Token = ``

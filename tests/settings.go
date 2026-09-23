package tests

import "go-final/pkg/config"

var Port = config.GetPortFromEnv()
var DBFile = "../scheduler.db"
var FullNextDate = false
var Search = false
var Token = ``

package tests

import (
	"fmt"
	"go-final/pkg/config"
)

var Port = config.GetPortFromEnv()
var DBFile = fmt.Sprintf("../%s", config.GetDBFromEnv())
var FullNextDate = true
var Search = true
var Token = getToken(config.GetPasswordFromEnv())

func getToken(pass string) string {
	token, err := config.GetToken(pass)
	if err != nil {
		return ""
	}
	fmt.Printf("Token: %s\n", token)
	return token
}

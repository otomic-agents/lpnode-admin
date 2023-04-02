package globalvar

import (
	"os"
	"strings"
)

var LpNodeHost = "lpnode-server"
var LpNodePort int = 9202
var SystemEnv string

func init() {
	SystemEnv = os.Getenv("DEPLOY_ENV")
	if SystemEnv == "" {
		SystemEnv = "dev"
	}
	SystemEnv = strings.ToLower(SystemEnv)
	if strings.Contains(SystemEnv, "prod") || strings.Contains(SystemEnv, "production") {
		SystemEnv = "prod"
	}
}

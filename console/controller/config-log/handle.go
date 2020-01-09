package config_log

import (
	"fmt"
	"net/http"
	"strings"

	module "github.com/wanghonggao007/goku-api-gateway/console/module/config-log"
)

//Handle handle
func Handle(prefix string) http.Handler {

	pre := strings.TrimSuffix(prefix, "/")
	serveMux := http.NewServeMux()
	consoleLogHandler := &LogHandler{
		name: module.ConsoleLog,
	}

	serveMux.Handle(fmt.Sprintf("%s/%s", pre, "console"), consoleLogHandler)
	nodeLogHandler := &LogHandler{
		name: module.NodeLog,
	}
	serveMux.Handle(fmt.Sprintf("%s/%s", pre, "node"), nodeLogHandler)

	accessLogHandler := &AccessLogHandler{}

	serveMux.Handle(fmt.Sprintf("%s/%s", pre, "access"), accessLogHandler)

	return serveMux

}

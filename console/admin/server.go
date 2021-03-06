package admin

import (
	"net/http"

	"github.com/wanghonggao007/goku-api-gateway/console/module/versionConfig"
)

//StartServer 开启admin服务
func StartServer(bind string) error {
	handler := router()
	versionConfig.InitVersionConfig()
	return http.ListenAndServe(bind, handler)
}

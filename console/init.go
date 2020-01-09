package console

import (
	_ "github.com/wanghonggao007/goku-api-gateway/console/updater/manager"
	graphite "github.com/wanghonggao007/goku-api-gateway/module/graphite/config"
	prometheus "github.com/wanghonggao007/goku-api-gateway/module/prometheus/config"
)

func moduleRegister() {
	prometheus.Register()
	//prometheus_pushgateway.Register()
	graphite.Register()
	//statsD.Register()
}

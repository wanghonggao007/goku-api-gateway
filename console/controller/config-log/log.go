package config_log

import (
	"fmt"
	"net/http"

	"github.com/wanghonggao007/goku-api-gateway/common/auto-form"
	"github.com/wanghonggao007/goku-api-gateway/console/controller"
	module "github.com/wanghonggao007/goku-api-gateway/console/module/config-log"
)

//LogHandler 日志处理器
type LogHandler struct {
	name string
}

func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := controller.CheckLogin(w, r, controller.OperationGatewayConfig, controller.OperationEDIT)
	if err != nil {
		return
	}

	switch r.Method {
	case http.MethodGet:
		{
			h.get(w, r)

		}

	case http.MethodPut:
		{
			h.set(w, r)

		}
	default:
		w.WriteHeader(404)
	}
}
func (h *LogHandler) get(w http.ResponseWriter, r *http.Request) {
	config, e := module.Get(h.name)
	if e = r.ParseForm(); e != nil {
		controller.WriteError(w, "270000", "data", "[Get]未知错误："+e.Error(), e)
		return
	}

	controller.WriteResultInfo(w,
		"data",
		"data",
		config)

}

func (h *LogHandler) set(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err = r.ParseForm(); err != nil {
		controller.WriteError(w, "260000", "data", "[param_check] Parse form body error | 解析form表单参数错误", err)
		return
	}

	param := new(module.PutParam)
	err = auto.SetValues(r.Form, param)
	if err != nil {
		controller.WriteError(w, "260000", "data", fmt.Sprintf("[param_check] %s", err.Error()), err)
		return
	}
	if param.Dir == "" {
		controller.WriteError(w, "260000", "data", "[param_check] inval dir", err)
		return
	}
	if param.File == "" {
		controller.WriteError(w, "260000", "data", "[param_check] inval file name", err)
		return
	}
	if param.Expire < module.ExpireDefault {
		controller.WriteError(w, "260000", "data", "[param_check] inval expire", nil)
		return
	}
	paramBase, err := param.Format()
	if err != nil {
		controller.WriteError(w, "260000", "data", fmt.Sprintf("[param_check] %s", err.Error()), err)
		return
	}
	err = module.Set(h.name, paramBase)
	if err != nil {
		controller.WriteError(w, "260000", "data", fmt.Sprintf("[mysql_error] %s", err.Error()), err)
		return
	}
	controller.WriteResultInfo(w,
		"data",
		"",
		nil)
}

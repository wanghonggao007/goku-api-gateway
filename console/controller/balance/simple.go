package balance

import (
	"net/http"

	"github.com/wanghonggao007/goku-api-gateway/console/controller"
	"github.com/wanghonggao007/goku-api-gateway/console/module/balance"
)

//GetSimpleList 获取简易列表
func GetSimpleList(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationLoadBalance, controller.OperationREAD)
	if e != nil {
		return
	}

	flag, result, err := balance.GetBalancNames()

	if !flag {
		controller.WriteError(httpResponse, "260000,", "balance", "[ERROR]Empty balance list!", err)
		return
	}
	controller.WriteResultInfo(httpResponse, "balance", "balanceNames", result)

	return

}

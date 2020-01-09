package balance

import (
	"fmt"
	"net/http"

	"github.com/wanghonggao007/goku-api-gateway/console/controller"
	"github.com/wanghonggao007/goku-api-gateway/console/module/balance"
)

//GetBalanceList 获取负载列表
func GetBalanceList(w http.ResponseWriter, r *http.Request) {
	_, e := controller.CheckLogin(w, r, controller.OperationLoadBalance, controller.OperationREAD)
	if e != nil {
		return
	}
	_ = r.ParseForm()

	keyword := r.FormValue("keyword")
	result, err := balance.Search(keyword)
	if err != nil {
		controller.WriteError(w, "260000", "balance", fmt.Sprintf("[ERROR] %s", err.Error()), err)
		return
	}
	controller.WriteResultInfo(w, "balance", "balanceList", result)

	return
}

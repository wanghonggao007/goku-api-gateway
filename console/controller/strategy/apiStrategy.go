package strategy

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/wanghonggao007/goku-api-gateway/console/controller"
	"github.com/wanghonggao007/goku-api-gateway/console/module/api"
	"github.com/wanghonggao007/goku-api-gateway/console/module/strategy"
)

//AddAPIToStrategy 将接口加入策略组
func AddAPIToStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationEDIT)
	if e != nil {
		return
	}

	strategyID := httpRequest.PostFormValue("strategyID")
	apiID := httpRequest.PostFormValue("apiID")
	apiArray := strings.Split(apiID, ",")

	flag, err := strategy.CheckStrategyIsExist(strategyID)
	if !flag {
		controller.WriteError(httpResponse,
			"240013",
			"apiStrategy",
			"[ERROR]The strategy does not exist!",
			err)
		return

	}
	flag, result, err := api.AddAPIToStrategy(apiArray, strategyID)
	if !flag {
		controller.WriteError(httpResponse,
			"240000",
			"apiStrategy",
			result,
			err)
		return

	}

	controller.WriteResultInfo(httpResponse, "apiStrategy", "", nil)

}

// ResetAPITargetOfStrategy 将接口加入策略组
func ResetAPITargetOfStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationEDIT)
	if e != nil {
		return
	}

	strategyID := httpRequest.PostFormValue("strategyID")
	target := httpRequest.PostFormValue("target")
	apiID := httpRequest.PostFormValue("apiID")
	aID, err := strconv.Atoi(apiID)
	if err != nil {
		controller.WriteError(httpResponse,
			"240013",
			"apiStrategy",
			"[ERROR]The strategy does not exist!",
			err)
		return

	}
	flag, err := strategy.CheckStrategyIsExist(strategyID)
	if !flag {
		controller.WriteError(httpResponse,
			"240013",
			"apiStrategy",
			"[ERROR]The strategy does not exist!",
			err)
		return

	}
	flag, result, err := api.SetTarget(aID, strategyID, target)
	if !flag {
		controller.WriteError(httpResponse,
			"240000",
			"apiStrategy",
			result,
			err)
		return

	}
	controller.WriteResultInfo(httpResponse, "apiStrategy", "", nil)

}

// BatchResetAPITargetOfStrategy 将接口加入策略组
func BatchResetAPITargetOfStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationEDIT)
	if e != nil {
		return
	}

	strategyID := httpRequest.PostFormValue("strategyID")
	target := httpRequest.PostFormValue("target")
	apiIDs := httpRequest.PostFormValue("apiIDs")
	ids := make([]int, 0)
	err := json.Unmarshal([]byte(apiIDs), &ids)
	if err != nil || len(ids) < 1 {
		controller.WriteError(httpResponse,
			"240004",
			"apiStrategy",
			"[ERROR]Illegal apiIDs!",
			err)
		return
	}
	flag, err := strategy.CheckStrategyIsExist(strategyID)
	if !flag {
		controller.WriteError(httpResponse,
			"240013",
			"apiStrategy",
			"[ERROR]The strategy does not exist!",
			err)
		return

	}
	flag, result, err := api.BatchSetTarget(ids, strategyID, target)
	if !flag {
		controller.WriteError(httpResponse,
			"240000",
			"apiStrategy",
			result,
			err)
		return

	}
	controller.WriteResultInfo(httpResponse, "apiStrategy", "", nil)

}

// GetAPIIDListFromStrategy 获取策略组接口ID列表
func GetAPIIDListFromStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationREAD)
	if e != nil {
		return
	}

	httpRequest.ParseForm()
	strategyID := httpRequest.Form.Get("strategyID")
	keyword := httpRequest.Form.Get("keyword")
	condition := httpRequest.Form.Get("condition")
	idStr := httpRequest.Form.Get("ids")
	balanceNames := httpRequest.Form.Get("balanceNames")

	op, err := strconv.Atoi(condition)
	if err != nil {
	}
	var ids []int
	var names []string
	if op > 0 {
		switch op {
		case 1, 2:
			{
				err := json.Unmarshal([]byte(balanceNames), &names)
				if err != nil || len(names) < 1 {
					controller.WriteError(httpResponse, "240001", "apiStrategy", "[ERROR]Illegal balanceNames!", err)
					return
				}
				break

			}
		case 3, 4:
			{
				err := json.Unmarshal([]byte(idStr), &ids)
				if err != nil || len(ids) < 1 {
					controller.WriteError(httpResponse, "240002", "apiStrategy", "[ERROR]Illegal ids!", err)
					return
				}
				break
			}
		default:
			{
				controller.WriteError(httpResponse, "240003", "apiStrategy", "[ERROR]Illegal condition!", err)
				return
			}
		}

	}

	_, result, err := api.GetAPIIDListFromStrategy(strategyID, keyword, op, ids, names)
	controller.WriteResultInfoWithPage(httpResponse, "apiStrategy", "apiIDList", result, &controller.PageInfo{
		ItemNum:  len(result),
		TotalNum: len(result),
	})
	return
}

// GetAPIListFromStrategy 获取策略组接口列表
func GetAPIListFromStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationREAD)
	if e != nil {
		return
	}

	httpRequest.ParseForm()
	strategyID := httpRequest.Form.Get("strategyID")
	keyword := httpRequest.Form.Get("keyword")
	condition := httpRequest.Form.Get("condition")
	idStr := httpRequest.Form.Get("ids")
	balanceNames := httpRequest.Form.Get("balanceNames")
	page := httpRequest.Form.Get("page")
	pageSize := httpRequest.Form.Get("pageSize")

	p, e := strconv.Atoi(page)
	if e != nil {
		p = 1
	}
	pSize, e := strconv.Atoi(pageSize)
	if e != nil {
		pSize = 15
	}

	op, err := strconv.Atoi(condition)
	if err != nil {

	}
	var ids []int
	var names []string
	if op > 0 {
		switch op {
		case 1, 2:
			{
				err := json.Unmarshal([]byte(balanceNames), &names)
				if err != nil || len(names) < 1 {
					controller.WriteError(httpResponse, "240001", "apiStrategy", "[ERROR]Illegal balanceNames!", err)
					return
				}
				break

			}
		case 3, 4:
			{
				err := json.Unmarshal([]byte(idStr), &ids)
				if err != nil || len(ids) < 1 {
					controller.WriteError(httpResponse, "240002", "apiStrategy", "[ERROR]Illegal ids!", err)
					return
				}
				break
			}
		default:
			{
				controller.WriteError(httpResponse, "240003", "apiStrategy", "[ERROR]Illegal condition!", err)
				return
			}
		}

	}

	_, result, count, err := api.GetAPIListFromStrategy(strategyID, keyword, op, p, pSize, ids, names)

	controller.WriteResultInfoWithPage(httpResponse, "apiStrategy", "apiList", result, &controller.PageInfo{
		ItemNum:  len(result),
		TotalNum: count,
		Page:     p,
		PageSize: pSize,
	})
	return
}

//CheckIsExistAPIInStrategy 检查插件是否添加进策略组
func CheckIsExistAPIInStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationEDIT)
	if e != nil {
		return
	}

	strategyID := httpRequest.PostFormValue("strategyID")
	apiID := httpRequest.PostFormValue("apiID")

	id, err := strconv.Atoi(apiID)
	if err != nil {
		controller.WriteError(httpResponse,
			"190001",
			"apiStrategy",
			"[ERROR]Illegal apiID",
			err)
		return

	}
	flag, _, err := api.CheckIsExistAPIInStrategy(id, strategyID)
	if !flag {
		controller.WriteError(httpResponse,
			"240000",
			"apiStrategy",
			"[ERROR]Can not find the api in strategy!",
			err)
		return

	}
	controller.WriteResultInfo(httpResponse, "apiStrategy", "", nil)

	return
}

// GetAPIIDListNotInStrategyByProject 获取未被该策略组绑定的接口ID列表(通过项目)
func GetAPIIDListNotInStrategyByProject(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationREAD)
	if e != nil {
		return
	}
	httpRequest.ParseForm()
	strategyID := httpRequest.Form.Get("strategyID")
	projectID := httpRequest.Form.Get("projectID")
	groupID := httpRequest.Form.Get("groupID")
	keyword := httpRequest.Form.Get("keyword")

	pjID, err := strconv.Atoi(projectID)
	if err != nil {
		controller.WriteError(httpResponse,
			"240008",
			"apiStrategy",
			"[ERROR]Illegal projectID!",
			err)
		return
	}
	gID, err := strconv.Atoi(groupID)
	if err != nil {
		if groupID != "" {
			controller.WriteError(httpResponse,
				"240009",
				"apiStrategy",
				"[ERROR]Illegal groupID!",
				err)
			return
		}
		gID = -1
	}
	_, result, _ := api.GetAPIIDListNotInStrategy(strategyID, pjID, gID, keyword)
	controller.WriteResultInfoWithPage(httpResponse, "apiStrategy", "apiIDList", result, &controller.PageInfo{
		ItemNum:  len(result),
		TotalNum: len(result),
	})
	return
}

//GetAPIListNotInStrategy 获取未被该策略组绑定的接口列表
func GetAPIListNotInStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationREAD)
	if e != nil {
		return
	}
	httpRequest.ParseForm()
	strategyID := httpRequest.Form.Get("strategyID")
	projectID := httpRequest.Form.Get("projectID")
	groupID := httpRequest.Form.Get("groupID")
	keyword := httpRequest.Form.Get("keyword")
	page := httpRequest.Form.Get("page")
	pageSize := httpRequest.Form.Get("pageSize")

	p, e := strconv.Atoi(page)
	if e != nil {
		p = 1
	}
	pSize, e := strconv.Atoi(pageSize)
	if e != nil {
		pSize = 15
	}

	pjID, err := strconv.Atoi(projectID)
	if err != nil {
		controller.WriteError(httpResponse,
			"240008",
			"apiStrategy",
			"[ERROR]Illegal projectID!",
			err)
		return
	}
	gID, err := strconv.Atoi(groupID)
	if err != nil {
		if groupID != "" {
			controller.WriteError(httpResponse,
				"240009",
				"apiStrategy",
				"[ERROR]Illegal groupID!",
				err)
			return
		}
		gID = -1
	}
	result := make([]map[string]interface{}, 0)
	_, result, count, err := api.GetAPIListNotInStrategy(strategyID, pjID, gID, p, pSize, keyword)
	controller.WriteResultInfoWithPage(httpResponse, "apiStrategy", "apiList", result, &controller.PageInfo{
		ItemNum:  len(result),
		TotalNum: count,
		Page:     p,
		PageSize: pSize,
	})
	return
}

//BatchDeleteAPIInStrategy 批量删除策略组接口
func BatchDeleteAPIInStrategy(httpResponse http.ResponseWriter, httpRequest *http.Request) {
	_, e := controller.CheckLogin(httpResponse, httpRequest, controller.OperationStrategy, controller.OperationEDIT)
	if e != nil {
		return
	}

	apiIDList := httpRequest.PostFormValue("apiIDList")
	strategyID := httpRequest.PostFormValue("strategyID")

	flag, result, err := api.BatchDeleteAPIInStrategy(apiIDList, strategyID)
	if !flag {

		controller.WriteError(httpResponse,
			"240000",
			"apiStrategy",
			result,
			err)
		return

	}
	controller.WriteResultInfo(httpResponse, "apiStrategy", "", nil)
}

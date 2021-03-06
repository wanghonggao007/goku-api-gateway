package gateway

import (
	v "github.com/wanghonggao007/goku-api-gateway/common/version"
	console_sqlite3 "github.com/wanghonggao007/goku-api-gateway/server/dao/console-sqlite3"
)

//BaseGatewayInfo 网关基本配置
type BaseGatewayInfo struct {
	NodeCount      int    `json:"nodeCount"`
	NodeStartCount int    `json:"nodeStartCount"`
	NodeStopCount  int    `json:"nodeStopCount"`
	ProjectCount   int    `json:"projectCount"`
	APICount       int    `json:"apiCount"`
	StrategyCount  int    `json:"strategyCount"`
	PluginCount    int    `json:"pluginCount"`
	ClusterCount   int    `json:"clusterCount"`
	Version        string `json:"version"`
}

//SystemInfo 系统配置
type SystemInfo struct {
	BaseInfo BaseGatewayInfo `json:"baseInfo"`
}

//GetGatewayConfig 获取网关配置
func GetGatewayConfig() (map[string]interface{}, error) {
	return console_sqlite3.GetGatewayConfig()
}

//EditGatewayBaseConfig 编辑网关基本配置
func EditGatewayBaseConfig(successCode string, nodeUpdatePeriod, monitorUpdatePeriod, timeout int) (bool, string, error) {
	flag, result, err := console_sqlite3.EditGatewayBaseConfig(successCode, nodeUpdatePeriod, monitorUpdatePeriod, timeout)
	return flag, result, err
}

//EditGatewayAlarmConfig 编辑网关告警配置
func EditGatewayAlarmConfig(apiAlertInfo, sender, senderPassword, smtpAddress string, alertStatus, smtpPort, smtpProtocol int) (bool, string, error) {
	flag, result, err := console_sqlite3.EditGatewayAlarmConfig(apiAlertInfo, sender, senderPassword, smtpAddress, alertStatus, smtpPort, smtpProtocol)

	return flag, result, err
}

//GetGatewayMonitorSummaryByPeriod 获取监控summary
func GetGatewayMonitorSummaryByPeriod() (bool, *SystemInfo, error) {

	nodeStartCount, nodeStopCount, projectCount, apiCount, strategyCount, e := console_sqlite3.GetGatewayInfo()
	if e != nil {
		return false, nil, e
	}
	info := new(SystemInfo)
	info.BaseInfo.PluginCount = console_sqlite3.GetPluginCount()
	info.BaseInfo.NodeCount = nodeStartCount + nodeStopCount
	info.BaseInfo.ProjectCount = projectCount
	info.BaseInfo.APICount = apiCount
	info.BaseInfo.StrategyCount = strategyCount
	info.BaseInfo.Version = v.Version
	info.BaseInfo.ClusterCount = console_sqlite3.GetClusterCount()

	return true, info, nil

}

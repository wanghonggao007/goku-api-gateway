package console_sqlite3

import (
	SQL "database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	database2 "github.com/wanghonggao007/goku-api-gateway/common/database"
)

var strategyPlugins = []string{"goku-oauth2_auth", "goku-rate_limiting", "goku-replay_attack_defender"}

//AddPluginToStrategy 新增策略组插件
func AddPluginToStrategy(pluginName, config, strategyID string) (bool, interface{}, error) {
	db := database2.GetConnection()
	// 查询接口是否添加该插件
	sql := "SELECT strategyID FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var id string
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&id)
	if err == nil {
		return false, "[ERROR]The strategy plugin is already exist", errors.New("[ERROR]The strategy plugin is already exist")
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	Tx, _ := db.Begin()
	result, err := Tx.Exec("INSERT INTO goku_conn_plugin_strategy (pluginName,pluginConfig,strategyID,createTime,updateTime,pluginStatus) VALUES (?,?,?,?,?,?);", pluginName, config, strategyID, now, now, 1)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to insert data", errors.New("[ERROR]Fail to insert data")
	}
	connID, err := result.LastInsertId()
	if err != nil {
		Tx.Rollback()
		panic(err)
		return false, "[ERROR]Fail to insert data", errors.New("[ERROR]Fail to insert data")
	}

	sql = "UPDATE goku_gateway_strategy SET updateTime = ? WHERE strategyID = ?;"
	_, err = Tx.Exec(sql, now, strategyID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to update data!", err
	}
	Tx.Commit()
	return true, connID, nil
}

//EditStrategyPluginConfig 新增策略组插件配置
func EditStrategyPluginConfig(pluginName, config, strategyID string) (bool, string, error) {
	db := database2.GetConnection()
	// 查询策略组是否添加该插件
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")
	updateTag := t.Format("20060102150405")
	sql := "SELECT strategyID FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var id string
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&id)
	if err != nil {
		return false, "[ERROR]The strategy plugin is not exist", errors.New("[ERROR]The strategy plugin is not exist")
	}
	Tx, _ := db.Begin()
	_, err = Tx.Exec("UPDATE goku_conn_plugin_strategy SET updateTag = ?,pluginConfig = ?,updateTime = ? WHERE pluginName = ? AND strategyID = ?;", updateTag, config, now, pluginName, strategyID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to update data", errors.New("[ERROR]Fail to update data")
	}

	sql = "UPDATE goku_gateway_strategy SET updateTime = ? WHERE strategyID = ?;"
	_, err = Tx.Exec(sql, now, strategyID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to update data!", err
	}
	Tx.Commit()
	return true, "", nil
}

// GetStrategyPluginList 获取策略插件列表
func GetStrategyPluginList(strategyID, keyword string, condition int) (bool, []map[string]interface{}, error) {
	db := database2.GetConnection()

	rule := make([]string, 0, 3)

	rule = append(rule, fmt.Sprintf("A.strategyID = '%s'", strategyID))
	if keyword != "" {
		searchRule := "(A.pluginName LIKE '%" + keyword + "%' OR B.pluginDesc LIKE '%" + keyword + "%')"
		rule = append(rule, searchRule)
	}
	if condition > 0 {
		rule = append(rule, fmt.Sprintf("CASE WHEN B.pluginStatus=0 THEN -1 ELSE A.pluginStatus = %d END", condition-1))
	}

	ruleStr := ""
	if len(rule) > 0 {
		ruleStr += "WHERE " + strings.Join(rule, " AND ")
	}
	sql := fmt.Sprintf(`SELECT A.connID,A.pluginName,A.pluginConfig,IFNULL(A.createTime,""),IFNULL(A.updateTime,""),B.pluginPriority,CASE WHEN B.pluginStatus=0 THEN -1 ELSE A.pluginStatus END as pluginStatus,IFNULL(B.pluginDesc,"") FROM goku_conn_plugin_strategy A INNER JOIN goku_plugin B ON B.pluginName = A.pluginName %s ORDER BY pluginStatus DESC,A.updateTime DESC;`, ruleStr)
	rows, err := db.Query(sql)
	if err != nil {
		return false, make([]map[string]interface{}, 0), err
	}
	defer rows.Close()
	pluginList := make([]map[string]interface{}, 0)
	//获取记录列
	for rows.Next() {
		var pluginPriority, pluginStatus, connID int
		var pluginName, pluginDesc, pluginConfig, createTime, updateTime string
		err = rows.Scan(&connID, &pluginName, &pluginConfig, &createTime, &updateTime, &pluginPriority, &pluginStatus, &pluginDesc)
		if err != nil {
		}
		pluginInfo := map[string]interface{}{
			"connID":         connID,
			"pluginName":     pluginName,
			"pluginConfig":   pluginConfig,
			"pluginPriority": pluginPriority,
			"pluginStatus":   pluginStatus,
			"createTime":     createTime,
			"updateTime":     updateTime,
			"pluginDesc":     pluginDesc,
		}
		pluginList = append(pluginList, pluginInfo)
	}
	return true, pluginList, nil
}

//GetStrategyPluginConfig 通过策略组ID获取配置信息
func GetStrategyPluginConfig(strategyID, pluginName string) (bool, string, error) {
	db := database2.GetConnection()
	sql := "SELECT pluginConfig FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var p string
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&p)
	if err != nil {
		if err == SQL.ErrNoRows {
			return false, "", errors.New("[ERROR]Can not find the plugin")
		}
		return false, "", err
	}
	return true, p, nil

}

//CheckPluginIsExistInStrategy 检查策略组是否绑定插件
func CheckPluginIsExistInStrategy(strategyID, pluginName string) (bool, error) {
	db := database2.GetConnection()
	sql := "SELECT strategyID FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var id string
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&id)
	if err != nil {
		return false, err
	}
	return true, err
}

//GetStrategyPluginStatus 检查策略组插件是否开启
func GetStrategyPluginStatus(strategyID, pluginName string) (bool, error) {
	db := database2.GetConnection()
	sql := "SELECT pluginStatus FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var pluginStatus int
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&pluginStatus)
	if err != nil {
		return false, err
	}
	if pluginStatus != 1 {
		return false, nil
	}
	return true, nil
}

//GetConnIDFromStrategyPlugin 获取Connid
func GetConnIDFromStrategyPlugin(pluginName, strategyID string) (bool, int, error) {
	db := database2.GetConnection()
	sql := "SELECT connID FROM goku_conn_plugin_strategy WHERE strategyID = ? AND pluginName = ?;"
	var connID int
	err := db.QueryRow(sql, strategyID, pluginName).Scan(&connID)
	if err != nil {
		return false, 0, err
	}
	return true, connID, nil
}

//BatchEditStrategyPluginStatus 批量修改策略组插件状态
func BatchEditStrategyPluginStatus(connIDList, strategyID string, pluginStatus int) (bool, string, error) {
	db := database2.GetConnection()
	t := time.Now()
	now := t.Format("2006-01-02 15:04:05")
	updateTag := t.Format("20060102150405")
	Tx, _ := db.Begin()
	sql := "UPDATE goku_conn_plugin_strategy SET updateTag = ?,pluginStatus = ?, updateTime = ? WHERE connID IN (" + connIDList + ");"
	_, err := Tx.Exec(sql, updateTag, pluginStatus, now)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to excute SQL statement!", err
	}

	sql = "UPDATE goku_gateway_strategy SET updateTime = ? WHERE strategyID = ?;"
	_, err = Tx.Exec(sql, now, strategyID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to update data!", err
	}
	Tx.Commit()
	return true, "", nil
}

//BatchDeleteStrategyPlugin 批量删除策略组插件
func BatchDeleteStrategyPlugin(connIDList, strategyID string) (bool, string, error) {
	db := database2.GetConnection()
	now := time.Now().Format("2006-01-02 15:04:05")
	Tx, _ := db.Begin()
	sql := "DELETE FROM goku_conn_plugin_strategy WHERE connID IN (" + connIDList + ");"
	_, err := Tx.Exec(sql)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to excute SQL statement!", err
	}

	sql = "UPDATE goku_gateway_strategy SET updateTime = ? WHERE strategyID = ?;"
	_, err = Tx.Exec(sql, now, strategyID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to update data!", err
	}
	Tx.Commit()
	return true, "", nil
}

//CheckStrategyPluginIsExistByConnIDList 通过connIDList判断插件是否存在
func CheckStrategyPluginIsExistByConnIDList(connIDList, pluginName string) (bool, error) {
	db := database2.GetConnection()
	sql := "SELECT pluginStatus FROM goku_conn_plugin_strategy WHERE connID IN (" + connIDList + ") AND pluginName = ?;"
	var pluginStatus int
	err := db.QueryRow(sql, pluginName).Scan(&pluginStatus)
	if err != nil {
		return false, err
	}
	return true, nil
}

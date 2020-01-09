package console_sqlite3

import (
	database2 "github.com/wanghonggao007/goku-api-gateway/common/database"
)

//AddStrategyGroup 新建策略组分组
func AddStrategyGroup(groupName string) (bool, interface{}, error) {
	db := database2.GetConnection()
	sql := "INSERT INTO goku_gateway_strategy_group (groupName) VALUES (?);"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, "[ERROR]Illegal SQL statement!", err
	}
	defer stmt.Close()
	r, err := stmt.Exec(groupName)
	if err != nil {
		return false, "[ERROR]Fail to insert data!", err
	}
	groupID, err := r.LastInsertId()
	if err != nil {
		return false, "[ERROR]Fail to insert data!", err
	}
	return true, groupID, nil
}

//EditStrategyGroup 修改策略组分组
func EditStrategyGroup(groupName string, groupID int) (bool, string, error) {
	db := database2.GetConnection()
	sql := "UPDATE goku_gateway_strategy_group SET groupName = ? WHERE groupID = ?;"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return false, "[ERROR]Illegal SQL statement!", err
	}
	defer stmt.Close()
	_, err = stmt.Exec(groupName, groupID)
	if err != nil {
		return false, "[ERROR]Fail to update data!", err
	}
	return true, "", nil
}

//DeleteStrategyGroup 删除策略组分组
func DeleteStrategyGroup(groupID int) (bool, string, error) {
	db := database2.GetConnection()
	// 查询该分组下所有策略组ID
	sql := "SELECT strategyID FROM goku_gateway_strategy WHERE groupID = ?;"
	rows, err := db.Query(sql, groupID)
	if err != nil {
		return false, "[ERROR]Illegal SQL statement!", err
	}
	strategyIDList := make([]string, 0)
	//延时关闭Rows
	defer rows.Close()
	//获取记录列

	for rows.Next() {
		var strategyID string
		err = rows.Scan(&strategyID)
		if err != nil {
			return false, "[ERROR]Fail to excute SQL statement!", err
		}
		strategyIDList = append(strategyIDList, strategyID)
	}
	Tx, _ := db.Begin()
	_, err = Tx.Exec("DELETE FROM goku_gateway_strategy_group WHERE groupID = ?;", groupID)
	if err != nil {
		Tx.Rollback()
		return false, "[ERROR]Fail to excute SQL statement!", err
	}
	if len(strategyIDList) > 0 {
		code := ""
		s := make([]interface{}, 0)
		for i := 0; i < len(strategyIDList); i++ {
			code += "?"
			if i < len(strategyIDList)-1 {
				code += ","
			}
			s = append(s, strategyIDList[i])
		}
		// 删除绑定的接口
		sql = "DELETE FROM goku_conn_strategy_api WHERE strategyID IN (" + code + ");"
		_, err = Tx.Exec(sql, s...)
		if err != nil {
			Tx.Rollback()
			return false, "[ERROR]Fail to excute SQL statement!", err
		}

		_, err = Tx.Exec("DELETE FROM goku_gateway_strategy WHERE strategyID IN ("+code+");", s...)
		if err != nil {
			Tx.Rollback()
			return false, "[ERROR]Fail to delete data!", err
		}
		_, err = Tx.Exec("DELETE FROM goku_conn_plugin_strategy WHERE strategyID IN ("+code+")", s...)
		if err != nil {
			Tx.Rollback()
			return false, "[ERROR]Fail to delete data!", err
		}

		_, err = Tx.Exec("DELETE FROM goku_conn_plugin_api WHERE strategyID IN ("+code+")", s...)
		if err != nil {
			Tx.Rollback()
			return false, "[ERROR]Fail to delete data!", err
		}
	}
	Tx.Commit()
	return true, "", nil
}

//GetStrategyGroupList 获取策略组分组列表
func GetStrategyGroupList() (bool, []map[string]interface{}, error) {
	db := database2.GetConnection()
	sql := "SELECT groupID,groupName,groupType FROM goku_gateway_strategy_group WHERE groupType = 0;"
	rows, err := db.Query(sql)
	if err != nil {
		return false, nil, err
	}
	defer rows.Close()
	//获取记录列

	groupList := make([]map[string]interface{}, 0)
	for rows.Next() {
		var groupID, groupType int
		var groupName string
		err = rows.Scan(&groupID, &groupName, &groupType)
		if err != nil {
			return false, nil, err
		}

		groupInfo := map[string]interface{}{
			"groupID":   groupID,
			"groupName": groupName,
		}
		groupList = append(groupList, groupInfo)
	}
	return true, groupList, nil
}

//CheckIsOpenGroup 判断是否是开放分组
func CheckIsOpenGroup(groupID int) bool {
	db := database2.GetConnection()
	var groupType int
	sql := "SELECT groupType FROM goku_gateway_strategy_group WHERE groupID = ?;"
	err := db.QueryRow(sql, groupID).Scan(&groupType)
	if err != nil {
		return false
	}
	if groupType == 1 {
		return true
	}
	return false
}

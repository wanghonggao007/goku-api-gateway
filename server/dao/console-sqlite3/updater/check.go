package updater

import (
	SQL "database/sql"
	"fmt"

	"github.com/wanghonggao007/goku-api-gateway/common/database"
)

//IsTableExist 检查table是否存在
func IsTableExist(name string) bool {
	db := database.GetConnection()
	selectType := ""
	sql := "SELECT type FROM sqlite_master WHERE `type` = 'table' AND `name` = ?"
	err := db.QueryRow(sql, name).Scan(&selectType)
	if err != nil {
		return false
	}
	return true
}

//IsColumnExist 检查列是否存在
func IsColumnExist(name string, column string) bool {
	db := database.GetConnection()
	sql := fmt.Sprintf("PRAGMA table_info(%s)", name)
	rows, err := db.Query(sql)
	if err != nil {
		panic(err)
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var cID, notNull, pk int
		var name, columnType string
		var dfltValue SQL.NullString
		err = rows.Scan(&cID, &name, &columnType, &notNull, &dfltValue, &pk)
		if err != nil {
			return false
		}
		if name == column {
			return true
		}
	}
	return false
}

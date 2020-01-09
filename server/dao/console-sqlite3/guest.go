package console_sqlite3

import (
	SQL "database/sql"

	"github.com/wanghonggao007/goku-api-gateway/common/database"
	"github.com/wanghonggao007/goku-api-gateway/utils"
)

//Login 登录
func Login(loginCall, loginPassword string) (bool, int) {
	db := database.GetConnection()
	var userID int
	err := db.QueryRow("SELECT userID FROM goku_admin WHERE loginCall = ? AND loginPassword = ?;", loginCall, loginPassword).Scan(&userID)
	if err != nil {
		return false, 0
	}
	return true, userID
}

//CheckLogin 检查用户是否登录
func CheckLogin(userToken string, userID int) bool {
	db := database.GetConnection()
	var loginPassword, loginCall string
	err := db.QueryRow("SELECT loginCall,loginPassword FROM goku_admin WHERE userID = ?;", userID).Scan(&loginCall, &loginPassword)
	if err != nil {
		return false
	}
	if utils.Md5(loginCall+loginPassword) == userToken {
		return true
	}
	return false
}

//Register 用户注册
func Register(loginCall, loginPassword string) bool {
	db := database.GetConnection()
	sql := "SELECT userID,loginPassword FROM goku_admin WHERE loginCall = ?;"
	password := ""
	userID := 0
	err := db.QueryRow(sql, loginCall).Scan(&userID, &password)
	if err != nil {
		if err == SQL.ErrNoRows {
			sql = "INSERT INTO goku_admin (loginPassword,loginCall) VALUES (?,?);"
		} else {
			return false
		}
	} else {
		if password != loginPassword {
			sql = "UPDATE goku_admin SET loginPassword = ? WHERE loginCall = ?;"
		} else {
			return true
		}
	}
	rows, err := db.Exec(sql, loginPassword, loginCall)
	if err != nil {
		return false
	}
	affectRow, _ := rows.RowsAffected()
	if affectRow > 0 {
		return true
	}
	return false
}

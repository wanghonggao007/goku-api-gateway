package console

import (
	"fmt"

	"github.com/wanghonggao007/goku-api-gateway/common/database"
	"github.com/wanghonggao007/goku-api-gateway/server/entity"
)

const (
	mysqlDriver   = "mysql"
	sqlite3Driver = "sqlite3"
)

//ClusterDatabaseConfig 集群数据库配置
type ClusterDatabaseConfig entity.ClusterDB

//GetDriver 获取驱动类型
func (c *ClusterDatabaseConfig) GetDriver() string {
	return c.Driver
}

//GetSource 获取连接字符串
func (c *ClusterDatabaseConfig) GetSource() string {

	switch c.Driver {
	case mysqlDriver:
		{
			return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", c.UserName, c.Password, c.Host, c.Port, c.Database)
		}
	case sqlite3Driver:
		{
			return c.Path
		}
	default:
		{
			return ""
		}
	}

}

//InitDatabase 初始化数据库
func InitDatabase() {
	def, err := getDefaultDatabase()
	if err != nil {
		panic(err)
	}
	c := ClusterDatabaseConfig(*def)
	fmt.Println(c)
	e := database.InitConnection(&c)
	if e != nil {
		//panic(e) //w---
		fmt.Println(e) //w+++
	}
}

func InitTable() error {
	return database.InitTable()
}

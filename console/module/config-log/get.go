package config_log

import (
	"fmt"

	"github.com/wanghonggao007/goku-api-gateway/common/auto-form"
	log "github.com/wanghonggao007/goku-api-gateway/goku-log"
	config_log "github.com/wanghonggao007/goku-api-gateway/server/dao/console-sqlite3/config-log"
)

//Get 获取普通日志配置
func Get(name string) (*LogConfig, error) {
	if _, has := logNames[name]; !has {
		return nil, fmt.Errorf("not has that log config of %s", name)
	}

	c := &LogConfig{}
	c.Levels = Levels
	c.Periods = Periods
	c.Expires = Expires
	config, e := config_log.Get(name)

	if e != nil || config == nil {
		auto.SetDefaults(c)
		c.Name = name
		c.File = name
		c.Level = log.ErrorLevel.String()
		c.Period = log.PeriodHour.String()
		c.Expire = ExpireDefault
	} else {
		c.Read(config)
	}

	return c, nil
}

//GetAccess 获取access配置
func GetAccess() (*AccessConfig, error) {
	config, e := config_log.Get(AccessLog)
	c := new(AccessConfig)
	c.Periods = Periods
	c.Expires = Expires
	if e != nil || config == nil {
		auto.SetDefaults(c)
		c.Name = AccessLog

		c.Period = log.PeriodHour.String()
		c.Expire = ExpireDefault
		c.InitFields()
	} else {
		c.Read(config)
	}
	return c, nil
}

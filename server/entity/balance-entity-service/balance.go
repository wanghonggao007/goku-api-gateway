package entity

import "github.com/wanghonggao007/goku-api-gateway/server/driver"

//Balance 负载
type Balance struct {
	Name          string
	ServiceName   string
	ServiceDriver string
	ServiceType   string
	AppName       string
	Static        string
	StaticCluster string
	Desc          string
	CreateTime    string
	UpdateTime    string
}

//Type 获取负载类型
func (e *Balance) Type() *Balance {

	if e != nil {

		d, has := driver.Get(e.ServiceDriver)
		if has {
			e.ServiceType = d.Type
		}
	}

	return e
}

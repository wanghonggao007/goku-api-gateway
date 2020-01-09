package config

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/wanghonggao007/goku-api-gateway/ksitigarbha"
)

//ModuleNameSpace 模块空间名称
const ModuleNameSpace = "diting.graphite"
const moduleName = "Graphite"
const desc = "API监控模块对接Graphite(udp by minute)"
const addressPattern = `^[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]$`
const content = `[
        {
            "type": "line",
            "label":"接入地址",
            "descript":"(仅支持UDP)",
            "items":[
                {
                    "type":"text",
                    "name":"accessAddress",
                    "placeholder":"",
                    "required":true,
                    "pattern":"` + addressPattern + `"
                }
            ]
        }
    ]`

var (
	mode           []ksitigarbha.Model
	addressMatcher *regexp.Regexp
)

func init() {
	json.Unmarshal([]byte(content), &mode)

	r, err := regexp.Compile(addressPattern)
	if err != nil {
		panic("init graphite module error:" + err.Error())
	}
	addressMatcher = r
}

//GraphiteModule 配置
type GraphiteModule struct {
}

//GraphiteConfig graphiteConfig
type GraphiteConfig struct {
	AccessAddress string `json:"accessAddress"`
}

//GetModel getModel
func (c *GraphiteModule) GetModel() []ksitigarbha.Model {
	return mode
}

//GetDesc getDesc
func (c *GraphiteModule) GetDesc() string {
	return desc
}

//GetName getName
func (c *GraphiteModule) GetName() string {
	return moduleName
}

//GetNameSpace getNameSpace
func (c *GraphiteModule) GetNameSpace() string {
	return ModuleNameSpace
}

//GetDefaultConfig getDefauleConfig
func (c *GraphiteModule) GetDefaultConfig() interface{} {
	return &GraphiteConfig{
		AccessAddress: "",
	}
}

//Encoder encoder
func (c *GraphiteModule) Encoder(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	if vm, ok := v.(*GraphiteConfig); ok {
		d, _ := json.Marshal(vm)
		return string(d), nil
	}

	return "", errors.New("illegal config")
}

//Decode decode
func Decode(config string) (*GraphiteConfig, error) {
	mc := new(GraphiteConfig)
	err := json.Unmarshal([]byte(config), &mc)
	if err != nil {
		return nil, err
	}

	match := addressMatcher.MatchString(mc.AccessAddress)

	if !match {
		return nil, errors.New("invalid accessAddress")
	}
	return mc, nil
}

//Decode decode
func (c *GraphiteModule) Decode(config string) (interface{}, error) {
	return Decode(config)
}

//Register 模板注册
func Register() {
	ksitigarbha.Register(moduleName, new(GraphiteModule))
}

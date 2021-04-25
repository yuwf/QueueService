package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

var theConf = &ServerConf{Port: 8001, MaxOnlineNum: 50000, LoginNumPreSec: 500}

// 服务器配置
type ServerConf struct {
	PProf          int `json:"pprof"`          // 开启性能监控
	Port           int `json:"port"`           // 监听客户端的端口
	MaxOnlineNum   int `json:"maxonlinenum"`   // 最大在线人数
	LoginNumPreSec int `json:"loginnumpresec"` // 每秒允许最大的登录人数
	Color          int `json:"color"`          // 终端输出是否带颜色信息 windows自动终端好像不支持 可以使用Cmder终端
}

// 读取配置文件
func (self *ServerConf) LoadConf(path string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("ServerConf Read Error")
		return false
	}
	dataJson := []byte(data)
	if err = json.Unmarshal(dataJson, self); err != nil {
		log.Error().Str("path", path).Err(err).Msg("ServerConf Json Error")
		return false
	}
	log.Info().Str("path", path).Msg("ServerConf Load Success")
	return true
}

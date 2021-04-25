package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/rs/zerolog/log"
)

var theConf = &ClientConf{ServerAddr: "127.0.0.1:8001", ClientNum: 1}

//服务器配置
type ClientConf struct {
	ServerAddr string `json:"serveraddr"` // 服务器地址 ip:port
	ClientNum  int    `json:"clientnum"`
	ShowQueue  int    `json:"showqueue"`
	Color      int    `json:"color"` // 终端输出是否带颜色信息 windows自动终端好像不支持 可以使用Cmder终端
}

//读取配置文件
func (self *ClientConf) LoadConf(path string) bool {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error().Str("path", path).Err(err).Msg("ClientConf Read Error")
		return false
	}
	dataJson := []byte(data)
	if err = json.Unmarshal(dataJson, self); err != nil {
		log.Error().Str("path", path).Err(err).Msg("ClientConf Json Error")
		return false
	}
	log.Info().Str("path", path).Msg("ClientConf Load Success")
	return true
}

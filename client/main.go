package main

import (
	"flag"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	//初始化命令行参数
	flag.Parse()

	//读取配置
	if !theConf.LoadConf("client.json") {
		//return false 配置读取不到暂时用默认值
	}

	//设置全局的zerolog对象的属性
	zerolog.SetGlobalLevel(zerolog.InfoLevel) //设置日志级别
	//log.Logger = log.With().Caller().Logger() //Add file and line number to log
	//log.Logger = log.With().Caller().Timestamp().Logger()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: theConf.Color == 0, TimeFormat: "2006-01-02 15:04:05.000"}) //控制台输出优化
	//通过钩子可以将日志输出文件

	if !theApp.Init() {
		return
	}

	theApp.Run()
}

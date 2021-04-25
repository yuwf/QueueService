package main

import (
	"github.com/rs/zerolog/log"
	"sync/atomic"
	"time"
)

var theApp = &ServerApp{}

//全局的程序主类
type ServerApp struct {
	MainChan chan string
}

func (self *ServerApp) Init() bool {
	log.Info().Msg("ServerApp Init")

	self.MainChan = make(chan string, 10)

	//队列模块初始化
	if !theQueue.Init() {
		return false
	}

	//服务模块初始化 开启网络监听
	if !theServer.Init() {
		return false
	}

	return true
}

func (self *ServerApp) Run() {
	log.Info().Msg("ServerApp Run")

	// 主命令
	for {
		select {
		case msg := <-self.MainChan:
			if msg == "quit" {
				goto quit
			}
		case <-time.After(time.Second):
			connnum := int(atomic.LoadInt32(&theServer.ConnNum))
			onlinenum := int(atomic.LoadInt32(&theServer.OnlineNum))
			lastconnnum := int(atomic.LoadInt32(&theServer.LastConnNum))
			lastloginnum := int(atomic.LoadInt32(&theServer.LastLoginNum))
			lastdisconnnum := int(atomic.LoadInt32(&theServer.LastDisConnNum))

			queueadd := int(atomic.LoadInt32(&theQueue.LastAddNum))
			queueremove := int(atomic.LoadInt32(&theQueue.LastRemoveNum))
			queuelogin := int(atomic.LoadInt32(&theQueue.LastLoginNum))

			log.Info().Int("queueadd", queueadd).Int("queueremove", queueremove).Int("queuelogin", queuelogin).Int("queuelen", theQueue.GetClientNum()).
				Int("online", onlinenum).Int("lastconn", lastconnnum).Int("lastlogin", lastloginnum).Int("lastdisconn", lastdisconnnum).Int("conn", connnum).Msg("ServerInfo")

			// 还原最近统计的数
			atomic.StoreInt32(&theServer.LastConnNum, 0)
			atomic.StoreInt32(&theServer.LastLoginNum, 0)
			atomic.StoreInt32(&theServer.LastDisConnNum, 0)

			atomic.StoreInt32(&theQueue.LastAddNum, 0)
			atomic.StoreInt32(&theQueue.LastRemoveNum, 0)
			atomic.StoreInt32(&theQueue.LastLoginNum, 0)
		}
	}

quit:
	self.Quit()
}

func (self *ServerApp) Quit() {
	log.Info().Msg("ServerApp Quit")

	theServer.Close()
	theQueue.Quit()

	close(self.MainChan)
}

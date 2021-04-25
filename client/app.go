package main

import (
	"QueueService/msg"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

var theApp = &ClientApp{clients: make([]*Client, 0)}

//全局的程序主类
type ClientApp struct {
	MainChan chan string

	clients []*Client
}

func (self *ClientApp) Init() bool {
	log.Info().Msg("ClientApp Init")

	self.MainChan = make(chan string, 10)

	log.Info().Str("Addr", theConf.ServerAddr).Int("ClientNum", theConf.ClientNum).Msg("ConfigInfo")

	self.clients = make([]*Client, theConf.ClientNum)
	for i := 0; i < theConf.ClientNum; i++ {
		self.clients[i] = NewClient()
		self.clients[i].Init()
	}

	return true
}

func (self *ClientApp) Run() {
	log.Info().Msg("ClientApp Run")

	// 主命令
	for {
		select {
		case msg := <-self.MainChan:
			if msg == "quit" {
				goto quit
			}
		case <-time.After(time.Second * 1):
			self.UpdateClient()
		}
	}

quit:
	self.Quit()
}

func (self *ClientApp) Quit() {
	log.Info().Msg("ClientApp Quit")

	for _, client := range self.clients {
		client.CloseServer()
	}

	close(self.MainChan)
}

func (self *ClientApp) UpdateClient() {
	invalid := 0
	loginqueue := 0
	logined := 0
	close := 0
	for _, client := range self.clients {
		c := client
		if c.LoginState == UserLoginState_Invalid {
			invalid++
		} else if c.LoginState == UserLoginState_LoginQueue {
			loginqueue++
		} else if c.LoginState == UserLoginState_Logined {
			logined++
		}

		c.Seq.Submit(func() {
			if c.LoginState == UserLoginState_LoginQueue {
				// 排队中 获取排队信息
				if theConf.ShowQueue != 0 {
					req := &msg.LoginQueueReq{}
					c.Conn.SendMsg(msg.MSGID_Req|msg.MSGID_LoginQueue, req)
				}
			} else if c.LoginState == UserLoginState_Logined {
				//req := &msg.LoginQueueReq{}
				//c.Conn.SendMsg(msg.MSGID_Req|msg.MSGID_LoginQueue, req)
				// 玩家10分钟内随机掉线
				if rand.Intn(10*60) == 100 {
					c.CloseServer() // 直接关闭 内部会重连
				}
			}
		})
	}
	log.Info().Int("invalid", invalid).Int("loginqueue", loginqueue).Int("logined", logined).Int("total", theConf.ClientNum).Int("close", close).Msg("ClientInfo")
}

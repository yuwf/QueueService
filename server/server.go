package main

import (
	"QueueService/utils"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/rs/zerolog/log"
)

var theServer = &Server{ConnMgr: new(sync.Map), ClientMgr: new(sync.Map)}

//服务器对象
type Server struct {
	*gnet.EventServer

	//所有的连接的客户端 [gnet.Conn:*Client]
	ConnMgr *sync.Map

	//验证为合法的客户端 [uid:*Client]
	ClientMgr *sync.Map

	// 下面的数据需要通过原子操作
	// 当前的连接数 和ConnMgr一致
	ConnNum int32
	// 在线人数(已经成功登录的人数) 当>=theConf.LoginNumPreSec时进入登录排队
	OnlineNum int32

	// theQueue中每秒刷新成0
	//最近一秒的连接人数
	LastConnNum int32
	// 最近一秒登录的人数 ，有登录用户+1，当>=theConf.LoginNumPreSec时进入登录排队
	LastLoginNum int32
	// 最近一秒的掉线人数
	LastDisConnNum int32
}

//服务器开启监听
func (self *Server) Init() bool {
	//开启监听 gnet.Serve会阻塞
	go func() {
		addr := fmt.Sprintf("tcp://:%d", theConf.Port)

		err := gnet.Serve(self, addr, gnet.WithMulticore(true), gnet.WithTCPKeepAlive(time.Minute*2), gnet.WithCodec(nil), gnet.WithReusePort(true))
		if err != nil {
			log.Error().Err(err).Msg("gnet.Serve err")
		}
	}()

	return true
}

func (self *Server) OnInitComplete(server gnet.Server) (action gnet.Action) {
	log.Info().Int("Port", theConf.Port).Msg("StartListen")
	return
}

func (self *Server) OnShutdown(server gnet.Server) {
	log.Debug().Int("Port", theConf.Port).Msg("StopListen")
}

func (self *Server) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	log.Debug().Str("ClientAddr", c.RemoteAddr().String()).Msg("OnOpened")

	//构造连接对象 和 一个客户端
	seq := &utils.Sequence{}
	conn := utils.NewGNetConn(c, seq)
	client := NewClient(conn, seq)
	self.ConnMgr.Store(c, client)
	atomic.AddInt32(&self.ConnNum, 1)
	atomic.AddInt32(&self.LastConnNum, 1)
	return
}

func (self *Server) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	log.Debug().Err(err).Str("ClientAddr", c.RemoteAddr().String()).Msg("OnClosed")

	client, ok := self.ConnMgr.Load(c)
	if ok {
		cli, _ := client.(*Client)
		cli.Seq.Submit(func() {
			cli.Close()
		})
		self.ConnMgr.Delete(c)
		atomic.AddInt32(&self.ConnNum, -1)
		atomic.AddInt32(&self.LastDisConnNum, 1)
	}

	return
}

func (self *Server) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {

	client, ok := self.ConnMgr.Load(c)
	if ok {
		client.(*Client).Conn.Recv(frame)
	}

	c.ResetBuffer()

	return
}

func (self *Server) Close() {
	self.ConnMgr.Range(func(key, value interface{}) bool {
		value.(*Client).Close()
		self.ConnMgr.Delete(key)
		return true
	})
}

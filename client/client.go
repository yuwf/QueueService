package main

import (
	"QueueService/msg"
	"QueueService/utils"
	"github.com/rs/zerolog/log"
	"net"
	"time"
)

type UserLoginState int

const (
	UserLoginState_Invalid    UserLoginState = 0 // 初始状态
	UserLoginState_LoginQueue UserLoginState = 1 // 登录排队中
	UserLoginState_Logined    UserLoginState = 2 // 登录成功
)

// 客户端对象
type Client struct {
	Conn *utils.Conn
	Seq  *utils.Sequence

	// 用户数据 用户数据的访问都需要通过Seq，保证数据的安全性
	Uid   string
	Token string
	// 用户状态
	LoginState UserLoginState
	// 用户切换到当前状态的时间
	LoginStateTime time.Time
}

func NewClient() *Client {
	client := &Client{
		Seq:            &utils.Sequence{},
		LoginState:     UserLoginState_Invalid,
		LoginStateTime: time.Now(),
	}

	return client
}

func (self *Client) Init() bool {

	self.Uid = utils.GetUserId()
	self.Seq.Submit(func() {
		self.ConnServer()
	})
	return true
}

func (self *Client) SetState(state UserLoginState) {
	self.LoginState = state
	self.LoginStateTime = time.Now()
}

func (self *Client) ConnServer() {
	self.CloseServer()

	tcpAddr, err := net.ResolveTCPAddr("tcp4", theConf.ServerAddr)
	if err != nil {
		log.Fatal().Err(err).Str("ServerAddr", theConf.ServerAddr).Msg("Server Addr Fail")
		return
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Error().Err(err).Msg("Server Connect Fail")
		self.Seq.Submit(func() {
			time.Sleep(time.Second * 5) // 延迟5秒调用
			self.ConnServer()
		})
		return
	}
	log.Debug().Str("ServerAddr", conn.RemoteAddr().String()).Msg("ConnServer Success")

	self.Conn = utils.NewConn(conn, self.Seq)
	//设置断开连接的回调
	self.Conn.CloseHandler = self.HandleDisConn
	//设置消息处理函数
	self.Conn.SetHandle(msg.MSGID_Ack|msg.MSGID_Login, self.HandleLoginAck, nil)
	self.Conn.SetHandle(msg.MSGID_Msg|msg.MSGID_Login, self.HandleLoginMsg, nil)
	self.Conn.SetHandle(msg.MSGID_Ack|msg.MSGID_LoginQueue, self.HandleLoginQueueAck, nil)
	self.Conn.SetHandle(msg.MSGID_Msg|msg.MSGID_LoginQueue, self.HandleLoginQueueMsg, nil)
	self.Conn.SetHandle(msg.MSGID_Ack|msg.MSGID_HeartBeat, self.HandleHeartBeatAck, nil)

	//开启连接对象的读写逻辑
	go self.Conn.StartRead()

	// 发送登录消息
	var loginReq = &msg.LoginReq{}
	loginReq.Uid = self.Uid
	self.Conn.SendMsg(msg.MSGID_Req|msg.MSGID_Login, loginReq)
}

func (self *Client) CloseServer() {
	if self.Conn != nil {
		log.Debug().Str("uid", self.Uid).Msg("CloseServer")
		self.Conn.Close()
		self.Conn = nil
	}
	self.SetState(UserLoginState_Invalid)
}

func (self *Client) HandleDisConn(conn *net.TCPConn) {
	// 重连
	self.SetState(UserLoginState_Invalid)

	self.Seq.Submit(func() {
		time.Sleep(time.Second) // 延迟1秒调用
		self.ConnServer()
	})
}

func (self *Client) HandleLoginAck(m interface{}, userData interface{}) {
	ack := m.(*msg.LoginAck)
	if ack.Result == 0 {
		self.SetState(UserLoginState_Logined)
		log.Debug().Str("uid", self.Uid).Str("Token", ack.Token).Msg("Login Success")
	} else if ack.Result == 1 {
		// uid重复了 重新生成一个 重新登录
		self.Uid = utils.GetUserId()
		var req = &msg.LoginReq{}
		req.Uid = self.Uid
		self.Conn.SendMsg(msg.MSGID_Req|msg.MSGID_Login, req)
		return
	} else if ack.Result == 2 {
		log.Debug().Str("uid", self.Uid).Msg("Begin Login Queue")
		self.SetState(UserLoginState_LoginQueue)
	} else {

	}
}

func (self *Client) HandleLoginMsg(m interface{}, userData interface{}) {
	msg := m.(*msg.LoginMsg)
	if msg.Result == 0 {
		log.Debug().Str("uid", self.Uid).Str("Token", msg.Token).Msg("Login Success")
		self.SetState(UserLoginState_Logined)
	} else {
		log.Error().Str("uid", self.Uid).Msg("Login fail")
	}
}

func (self *Client) HandleLoginQueueAck(m interface{}, userData interface{}) {
	ack := m.(*msg.LoginQueueAck)
	if ack.Result == 0 {
		log.Info().Int32("queuepos", ack.Pos).Int32("queuelen", ack.Num).Str("UId", self.Uid).Msg("LoginQueue")
	}
}

func (self *Client) HandleLoginQueueMsg(m interface{}, userData interface{}) {
	msg := m.(*msg.LoginQueueMsg)

	log.Debug().Int32("queuepos", msg.Pos).Int32("queuelen", msg.Num).Str("UId", self.Uid).Msg("LoginQueue")
}

func (self *Client) HandleHeartBeatAck(m interface{}, userData interface{}) {
	//heatBeatAck := m.(*msg.HeartBeatAck)
}

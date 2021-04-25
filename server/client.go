package main

import (
	"QueueService/msg"
	"QueueService/utils"
	"sync/atomic"
)

type UserLoginState int

const (
	UserLoginState_Invalid    UserLoginState = 0 // 初始状态
	UserLoginState_LoginQueue UserLoginState = 1 // 登录排队中
	UserLoginState_Logined    UserLoginState = 2 // 登录成功
)

// 客户端对象 处理客户端逻辑
//
type Client struct {
	// 顺序处理 Client所有的调用或者数据访问都应通过Seq
	Seq *utils.Sequence

	// 连接对象
	Conn *utils.GNetConn

	// 用户数据 用户数据的访问都需要通过Seq，保证数据的安全性
	// 用户状态
	LoginState UserLoginState
	// 用户ID
	UId string
}

func NewClient(conn *utils.GNetConn, s *utils.Sequence) *Client {
	client := &Client{
		Seq:        s,
		Conn:       conn,
		LoginState: UserLoginState_Invalid,
	}
	//设置消息处理函数
	conn.SetHandle(msg.MSGID_Req|msg.MSGID_Login, client.HandleLoginReq, nil)
	conn.SetHandle(msg.MSGID_Req|msg.MSGID_LoginQueue, client.HandleLoginQueueReq, nil)
	conn.SetHandle(msg.MSGID_Req|msg.MSGID_HeartBeat, client.HandleHeartBeatReq, nil)
	return client
}

func (self *Client) Close() {
	if self.LoginState == UserLoginState_Invalid {

	} else if self.LoginState == UserLoginState_LoginQueue {
		theQueue.RemoveClient(self.UId)
	} else if self.LoginState == UserLoginState_Logined {
		atomic.AddInt32(&theServer.OnlineNum, -1) // 在线人数-1
	}
	if len(self.UId) > 0 {
		theServer.ClientMgr.Delete(self.UId)
	}
}

func (self *Client) HandleLoginReq(m interface{}, userData interface{}) {
	req := m.(*msg.LoginReq)

	// 验证用户的合法性
	// todo ..
	// 认为用户为合法用户

	if self.LoginState == UserLoginState_Invalid {
		// OK
	} else if self.LoginState == UserLoginState_LoginQueue {
		// 用户排队中 直接给用户发送排队位置
		smsg := &msg.LoginQueueMsg{}
		pos, num := theQueue.GetClientPos(self.UId)
		smsg.Pos = int32(pos)
		smsg.Num = int32(num)
		self.Conn.SendMsg(msg.MSGID_Msg|msg.MSGID_LoginQueue, smsg)
		return
	} else if self.LoginState == UserLoginState_Logined {
		// 用户已经是登录状态了 先忽略
		return
	}

	// 判断用户是否已经存在
	_, ok := theServer.ClientMgr.Load(req.Uid)
	if ok {
		// 暂时直接发送登录失败
		ack := &msg.LoginAck{}
		ack.Result = 1
		self.Conn.SendMsg(msg.MSGID_Ack|msg.MSGID_Login, ack)
		return
	}

	// 保存信息
	self.UId = req.Uid
	theServer.ClientMgr.Store(req.Uid, self)

	// 判断用户是直接登录成功还是进入登录排队中
	onlinenum := int(atomic.LoadInt32(&theServer.OnlineNum))
	lastloginnum := int(atomic.LoadInt32(&theServer.LastLoginNum))
	queuelen := theQueue.GetClientNum()
	if onlinenum < theConf.MaxOnlineNum && lastloginnum < theConf.LoginNumPreSec && queuelen == 0 {
		// 玩家直接登录成功
		self.LoginState = UserLoginState_Logined
		atomic.AddInt32(&theServer.OnlineNum, 1)    // 在线人数+1
		atomic.AddInt32(&theServer.LastLoginNum, 1) // 登录人数+1

		sack := &msg.LoginAck{}
		sack.Result = 0
		sack.Token = utils.GetToken()
		self.Conn.SendMsg(msg.MSGID_Ack|msg.MSGID_Login, sack)
	} else {
		// 用户进入排队 并发给用户发送排队信息
		self.LoginState = UserLoginState_LoginQueue
		theQueue.AddClient(self.UId, self.Seq)

		var sack = &msg.LoginAck{}
		sack.Result = 2
		self.Conn.SendMsg(msg.MSGID_Ack|msg.MSGID_Login, sack)

		smsg := &msg.LoginQueueMsg{}
		pos, num := theQueue.GetClientPos(self.UId)
		smsg.Pos = int32(pos)
		smsg.Num = int32(num)
		self.Conn.SendMsg(msg.MSGID_Msg|msg.MSGID_LoginQueue, smsg)
	}
}

func (self *Client) HandleLoginQueueReq(m interface{}, userData interface{}) {
	ack := &msg.LoginQueueAck{}
	if self.LoginState == UserLoginState_LoginQueue {
		// 获取用户排队位置
		pos, num := theQueue.GetClientPos(self.UId)
		ack.Result = 0
		ack.Pos = int32(pos)
		ack.Num = int32(num)
	} else {
		// 用户没有排队
		ack.Result = 1
		ack.Pos = 0
		ack.Num = 0
	}

	self.Conn.SendMsg(msg.MSGID_Ack|msg.MSGID_LoginQueue, ack)
}

func (self *Client) HandleHeartBeatReq(m interface{}, userData interface{}) {
	//req := m.(*msg.HeartBeatReq)
	var ack = &msg.HeartBeatAck{}
	self.Conn.SendMsg(msg.MSGID_Ack|msg.MSGID_HeartBeat, ack)
}

// 完成排队 登录成功
func (self *Client) FinishQueue() {
	// 玩家直接登录成功
	self.LoginState = UserLoginState_Logined
	atomic.AddInt32(&theServer.OnlineNum, 1)    // 在线人数+
	atomic.AddInt32(&theServer.LastLoginNum, 1) // 最近登录数+

	// 发送登陆消息
	smsg := &msg.LoginMsg{}
	smsg.Result = 0
	smsg.Token = utils.GetToken()
	self.Conn.SendMsg(msg.MSGID_Msg|msg.MSGID_Login, smsg)
}

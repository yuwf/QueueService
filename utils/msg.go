package utils

import (
	"QueueService/msg"
	"fmt"
	"reflect"
	"unsafe"

	proto "github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
)

func init() {
	TheMsgMgr.Register(msg.MSGID_Req|msg.MSGID_Login, &msg.LoginReq{})
	TheMsgMgr.Register(msg.MSGID_Ack|msg.MSGID_Login, &msg.LoginAck{})
	TheMsgMgr.Register(msg.MSGID_Msg|msg.MSGID_Login, &msg.LoginMsg{})
	TheMsgMgr.Register(msg.MSGID_Req|msg.MSGID_LoginQueue, &msg.LoginQueueReq{})
	TheMsgMgr.Register(msg.MSGID_Ack|msg.MSGID_LoginQueue, &msg.LoginQueueAck{})
	TheMsgMgr.Register(msg.MSGID_Msg|msg.MSGID_LoginQueue, &msg.LoginQueueMsg{})
	TheMsgMgr.Register(msg.MSGID_Req|msg.MSGID_HeartBeat, &msg.HeartBeatReq{})
	TheMsgMgr.Register(msg.MSGID_Ack|msg.MSGID_HeartBeat, &msg.HeartBeatAck{})
}

// 消息头
type MsgHead struct {
	index    uint32 // 消息索引
	id       uint32 // 消息ID
	flag     uint16 // 消息标示
	ext      uint16 // 扩展使用
	datasize uint16 // 数据大小 不包括头
	checknum uint16 // 消息头校检和
}

const MsgHeadSize = int(unsafe.Sizeof(MsgHead{}))
const MsgHeadFlag uint16 = uint16(9527)

var TheMsgMgr = &MsgManager{make(map[msg.MSGID]*MsgInfo)}

type MsgHandler func(msg interface{}, userData interface{})

type MsgInfo struct {
	msgId   msg.MSGID    // 消息ID
	msgType reflect.Type // 消息类型
}

type MsgHandlerInfo struct {
	msgId    msg.MSGID // 消息ID
	handler  MsgHandler
	userData interface{}
}

//消息解析映射
type MsgManager struct {
	msgInfo map[msg.MSGID]*MsgInfo //msg.MSGID:MsgInfo 解析消息使用 程序启动时注册 不用加锁
}

// 注册消息号、消息类型和处理函数
func (self *MsgManager) Register(msgId msg.MSGID, msg proto.Message) {
	msgType := reflect.TypeOf(msg)
	if msgType == nil || msgType.Kind() != reflect.Ptr {
		log.Fatal().Msg("protobuf message pointer required")
		return
	}
	if _, ok := self.msgInfo[msgId]; ok {
		log.Error().Uint32("msgId", uint32(msgId)).Msg("MsgInfo already exist")
	}

	i := &MsgInfo{msgId, msgType}
	self.msgInfo[msgId] = i
}

// 解析消息
func (self *MsgManager) Unmarshal(msgId msg.MSGID, data []byte) (interface{}, error) {
	if _, ok := self.msgInfo[msgId]; !ok {
		log.Error().Uint32("msgId", uint32(msgId)).Msg("MsgInfo not exist")
		return nil, fmt.Errorf("MsgInfo not exist, msgId:%v", msgId)
	}

	// msg
	i := self.msgInfo[msgId]
	msg := reflect.New(i.msgType.Elem()).Interface()
	err := proto.UnmarshalMerge(data, msg.(proto.Message))
	if err != nil {
		log.Error().Err(err).Uint32("msgId", uint32(msgId)).Msg("MsgInfo Unmarshal err")
		return nil, fmt.Errorf("sgInfo Unmarshal err, msgId:%v", msgId)
	}

	return msg, nil
}

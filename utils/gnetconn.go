package utils

import (
	"QueueService/msg"
	"bytes"
	"encoding/binary"
	proto "github.com/golang/protobuf/proto"
	"github.com/panjf2000/gnet"
	"github.com/rs/zerolog/log"
	"strconv"
	"sync/atomic"
)

//gnet连接对象
type GNetConn struct {
	gconn     gnet.Conn
	Seq       *Sequence
	msgHandle map[msg.MSGID]*MsgHandlerInfo

	// 读取缓冲buff
	readBuff *bytes.Buffer
	// 标记当前是否已经读取了头
	readhead bool
	head     MsgHead

	// 消息索引
	recvindex uint32
	sendindex uint32 // 需要原子操作
}

func NewGNetConn(c gnet.Conn, s *Sequence) *GNetConn {
	conn := &GNetConn{
		gconn:     c,
		Seq:       s,
		msgHandle: make(map[msg.MSGID]*MsgHandlerInfo),
		readBuff:  new(bytes.Buffer),
		readhead:  false,
	}
	return conn
}

func (self *GNetConn) SetHandle(msgId msg.MSGID, handler MsgHandler, userData interface{}) {
	handle := &MsgHandlerInfo{msgId, handler, userData}
	self.msgHandle[msgId] = handle
}

//收到数据时调用
func (self *GNetConn) Recv(buf []byte) {

	//buf.Grow(n)
	self.readBuff.Write(buf)

	// 读取消息
	for {
		if !self.readhead {
			if self.readBuff.Len() >= MsgHeadSize {
				binary.Read(self.readBuff, binary.BigEndian, &self.head.flag)
				binary.Read(self.readBuff, binary.BigEndian, &self.head.index)
				binary.Read(self.readBuff, binary.BigEndian, &self.head.id)
				binary.Read(self.readBuff, binary.BigEndian, &self.head.ext)
				binary.Read(self.readBuff, binary.BigEndian, &self.head.datasize)
				binary.Read(self.readBuff, binary.BigEndian, &self.head.checknum)
				if self.head.flag != MsgHeadFlag {
					log.Error().Msg("msghead flag err")
					self.gconn.Close()
					return
				}
				self.readhead = true
			} else {
				break //还有粘包数据
			}
		} else {
			if self.readBuff.Len() >= int(self.head.datasize) {
				body := make([]byte, self.head.datasize)
				_, err := self.readBuff.Read(body)
				if err != nil {
					log.Info().Err(err).Msg("buff Read err")
					return
				}
				self.readhead = false
				self.recvindex++
				//解析数据
				m, err := TheMsgMgr.Unmarshal(msg.MSGID(self.head.id), body)
				if err == nil {
					log.Debug().Str("MsgId", strconv.FormatInt(int64(self.head.id), 16)).Str("Msg", m.(proto.Message).String()).Msg("RecvMsg")
					handle, ok := self.msgHandle[msg.MSGID(self.head.id)]
					if ok && self.Seq != nil {
						self.Seq.Submit(func() {
							HandlePanic()
							handle.handler(m, handle.userData)
						})
					}
				}
			} else {
				break //还有粘包数据
			}
		}
	}
}

func (self *GNetConn) SendMsg(msgId msg.MSGID, msg interface{}) bool {
	m, ok := msg.(proto.Message)
	if !ok {
		log.Error().Str("MsgId", strconv.FormatInt(int64(msgId), 16)).Msg("SendMsg protobuf message pointer required")
		return false
	}

	data, err := proto.Marshal(m)
	if err != nil {
		log.Error().Err(err).Msg("proto Marshal err")
		return false
	}

	// 写入消息头
	index := atomic.AddUint32(&self.sendindex, 1)
	var head MsgHead
	head.flag = MsgHeadFlag
	head.index = index
	head.id = uint32(msgId)
	head.ext = 0
	head.datasize = uint16(len(data))
	head.checknum = 0

	var buf = new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, head.flag)
	binary.Write(buf, binary.BigEndian, head.index)
	binary.Write(buf, binary.BigEndian, head.id)
	binary.Write(buf, binary.BigEndian, head.ext)
	binary.Write(buf, binary.BigEndian, head.datasize)
	binary.Write(buf, binary.BigEndian, head.checknum)

	// 写入消息实体
	_, err = buf.Write(data)
	if err != nil {
		log.Error().Err(err).Msg("binary Write err")
		return false
	}

	// 写入消息发送的channel中
	log.Debug().Str("MsgId", strconv.FormatInt(int64(msgId), 16)).Int("MsgSize", buf.Len()).Str("Msg", m.String()).Msg("SendMsg")
	self.gconn.AsyncWrite(buf.Bytes())
	return true
}

package utils

import (
	"QueueService/msg"
	"bytes"
	"encoding/binary"
	proto "github.com/golang/protobuf/proto"
	"github.com/rs/zerolog/log"
	"net"
	"reflect"
	"strconv"
	"sync/atomic"
)

//连接对象的回调函数
type ConnCloseHandler func(conn *net.TCPConn)

//连接对象
type Conn struct {
	CloseHandler ConnCloseHandler

	tcpConn   *net.TCPConn
	Seq       *Sequence
	msgHandle map[msg.MSGID]*MsgHandlerInfo

	// 读取缓冲buff
	readBuff *bytes.Buffer
	// 标记当前是否已经读取了头
	readhead bool
	head     MsgHead

	// 消息发送func执行序列
	writeSeq Sequence

	// 关闭控制
	closed bool // 防止重复调用

	// 消息索引
	recvindex uint32
	sendindex uint32 // 需要原子操作
}

func NewConn(c *net.TCPConn, s *Sequence) *Conn {
	conn := &Conn{
		tcpConn:   c,
		Seq:       s,
		msgHandle: make(map[msg.MSGID]*MsgHandlerInfo),
		closed:    false,
	}
	return conn
}

func (self *Conn) SetHandle(msgId msg.MSGID, handler MsgHandler, userData interface{}) {
	handle := &MsgHandlerInfo{msgId, handler, userData}
	self.msgHandle[msgId] = handle
}

func (self *Conn) StartRead() {
	defer HandlePanic()
	// 退出关闭
	defer self.Close()

	// 读取消息
	var head MsgHead
	var readhead bool = false // 标记当前是否已经读取了头
	var buf = new(bytes.Buffer)
	for {
		readbuf := make([]byte, 1024)
		n, err := self.tcpConn.Read(readbuf)
		if err != nil {
			log.Debug().Err(err).Msg("TCPConn Read err")
			break
		}
		//buf.Grow(n)
		buf.Write(readbuf[0:n])

		for {
			if !readhead {
				if buf.Len() >= MsgHeadSize {
					binary.Read(buf, binary.BigEndian, &head.flag)
					binary.Read(buf, binary.BigEndian, &head.index)
					binary.Read(buf, binary.BigEndian, &head.id)
					binary.Read(buf, binary.BigEndian, &head.ext)
					binary.Read(buf, binary.BigEndian, &head.datasize)
					binary.Read(buf, binary.BigEndian, &head.checknum)
					if head.flag != MsgHeadFlag {
						log.Info().Msg("msghead flag err")
						return
					}
					readhead = true
				} else {
					break //继续读取
				}
			} else {
				if buf.Len() >= int(head.datasize) {
					body := make([]byte, head.datasize)
					_, err := buf.Read(body)
					if err != nil {
						log.Info().Err(err).Msg("buff Read err")
						return
					}
					readhead = false
					self.recvindex++
					//
					m, err := TheMsgMgr.Unmarshal(msg.MSGID(head.id), body)
					if err == nil {
						log.Debug().Str("MsgId", strconv.FormatInt(int64(head.id), 16)).Int("MsgSize", len(body)+MsgHeadSize).Int("LeftBuf", buf.Len()).Str("Msg", m.(proto.Message).String()).Msg("RecvMsg")
						handle, ok := self.msgHandle[msg.MSGID(head.id)]
						if ok && self.Seq != nil {
							self.Seq.Submit(func() {
								HandlePanic()
								handle.handler(m, handle.userData)
							})
						}
					}
				} else {
					break // 继续读取
				}
			}
		}
	}
}

func (self *Conn) SendMsg(msgId msg.MSGID, msg interface{}) bool {
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

	log.Debug().Str("MsgId", strconv.FormatInt(int64(msgId), 16)).Int("MsgSize", buf.Len()).Str("Msg", m.String()).Str("MsgName", reflect.TypeOf(msg).Name()).Msg("SendMsg")
	// 发送消息
	self.writeSeq.Submit(func() {
		self.tcpConn.Write(buf.Bytes())
	})
	return true
}

func (self *Conn) Close() {
	if self.closed {
		return
	}
	self.closed = true

	log.Debug().Str("ClientAddr", self.tcpConn.RemoteAddr().String()).Msg("Client Close")

	if self.CloseHandler != nil {
		self.CloseHandler(self.tcpConn)
	}

	self.tcpConn.Close()
}

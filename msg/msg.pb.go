// Code generated by protoc-gen-go. DO NOT EDIT.
// source: msg.proto

package msg

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type MSGID int32

const (
	MSGID_FLAG       MSGID = 0
	MSGID_Req        MSGID = 4096
	MSGID_Ack        MSGID = 8192
	MSGID_Msg        MSGID = 16384
	MSGID_Login      MSGID = 16
	MSGID_LoginQueue MSGID = 18
	MSGID_HeartBeat  MSGID = 32
)

var MSGID_name = map[int32]string{
	0:     "FLAG",
	4096:  "Req",
	8192:  "Ack",
	16384: "Msg",
	16:    "Login",
	18:    "LoginQueue",
	32:    "HeartBeat",
}

var MSGID_value = map[string]int32{
	"FLAG":       0,
	"Req":        4096,
	"Ack":        8192,
	"Msg":        16384,
	"Login":      16,
	"LoginQueue": 18,
	"HeartBeat":  32,
}

func (x MSGID) String() string {
	return proto.EnumName(MSGID_name, int32(x))
}

func (MSGID) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

// 用户登录
type LoginReq struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{0}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

// 用户登录回复
type LoginAck struct {
	Result               int32    `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginAck) Reset()         { *m = LoginAck{} }
func (m *LoginAck) String() string { return proto.CompactTextString(m) }
func (*LoginAck) ProtoMessage()    {}
func (*LoginAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{1}
}

func (m *LoginAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginAck.Unmarshal(m, b)
}
func (m *LoginAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginAck.Marshal(b, m, deterministic)
}
func (m *LoginAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginAck.Merge(m, src)
}
func (m *LoginAck) XXX_Size() int {
	return xxx_messageInfo_LoginAck.Size(m)
}
func (m *LoginAck) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginAck.DiscardUnknown(m)
}

var xxx_messageInfo_LoginAck proto.InternalMessageInfo

func (m *LoginAck) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *LoginAck) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

// 排完对之后时发送
type LoginMsg struct {
	Result               int32    `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginMsg) Reset()         { *m = LoginMsg{} }
func (m *LoginMsg) String() string { return proto.CompactTextString(m) }
func (*LoginMsg) ProtoMessage()    {}
func (*LoginMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{2}
}

func (m *LoginMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginMsg.Unmarshal(m, b)
}
func (m *LoginMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginMsg.Marshal(b, m, deterministic)
}
func (m *LoginMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginMsg.Merge(m, src)
}
func (m *LoginMsg) XXX_Size() int {
	return xxx_messageInfo_LoginMsg.Size(m)
}
func (m *LoginMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginMsg.DiscardUnknown(m)
}

var xxx_messageInfo_LoginMsg proto.InternalMessageInfo

func (m *LoginMsg) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *LoginMsg) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

// 请求排队信息
type LoginQueueReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginQueueReq) Reset()         { *m = LoginQueueReq{} }
func (m *LoginQueueReq) String() string { return proto.CompactTextString(m) }
func (*LoginQueueReq) ProtoMessage()    {}
func (*LoginQueueReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{3}
}

func (m *LoginQueueReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginQueueReq.Unmarshal(m, b)
}
func (m *LoginQueueReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginQueueReq.Marshal(b, m, deterministic)
}
func (m *LoginQueueReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginQueueReq.Merge(m, src)
}
func (m *LoginQueueReq) XXX_Size() int {
	return xxx_messageInfo_LoginQueueReq.Size(m)
}
func (m *LoginQueueReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginQueueReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginQueueReq proto.InternalMessageInfo

type LoginQueueAck struct {
	Result               int32    `protobuf:"varint,1,opt,name=result,proto3" json:"result,omitempty"`
	Pos                  int32    `protobuf:"varint,2,opt,name=pos,proto3" json:"pos,omitempty"`
	Num                  int32    `protobuf:"varint,3,opt,name=num,proto3" json:"num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginQueueAck) Reset()         { *m = LoginQueueAck{} }
func (m *LoginQueueAck) String() string { return proto.CompactTextString(m) }
func (*LoginQueueAck) ProtoMessage()    {}
func (*LoginQueueAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{4}
}

func (m *LoginQueueAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginQueueAck.Unmarshal(m, b)
}
func (m *LoginQueueAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginQueueAck.Marshal(b, m, deterministic)
}
func (m *LoginQueueAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginQueueAck.Merge(m, src)
}
func (m *LoginQueueAck) XXX_Size() int {
	return xxx_messageInfo_LoginQueueAck.Size(m)
}
func (m *LoginQueueAck) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginQueueAck.DiscardUnknown(m)
}

var xxx_messageInfo_LoginQueueAck proto.InternalMessageInfo

func (m *LoginQueueAck) GetResult() int32 {
	if m != nil {
		return m.Result
	}
	return 0
}

func (m *LoginQueueAck) GetPos() int32 {
	if m != nil {
		return m.Pos
	}
	return 0
}

func (m *LoginQueueAck) GetNum() int32 {
	if m != nil {
		return m.Num
	}
	return 0
}

// 用户排队信息 进入排队或者排队信息变化发送
type LoginQueueMsg struct {
	Pos                  int32    `protobuf:"varint,1,opt,name=pos,proto3" json:"pos,omitempty"`
	Num                  int32    `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginQueueMsg) Reset()         { *m = LoginQueueMsg{} }
func (m *LoginQueueMsg) String() string { return proto.CompactTextString(m) }
func (*LoginQueueMsg) ProtoMessage()    {}
func (*LoginQueueMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{5}
}

func (m *LoginQueueMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginQueueMsg.Unmarshal(m, b)
}
func (m *LoginQueueMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginQueueMsg.Marshal(b, m, deterministic)
}
func (m *LoginQueueMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginQueueMsg.Merge(m, src)
}
func (m *LoginQueueMsg) XXX_Size() int {
	return xxx_messageInfo_LoginQueueMsg.Size(m)
}
func (m *LoginQueueMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginQueueMsg.DiscardUnknown(m)
}

var xxx_messageInfo_LoginQueueMsg proto.InternalMessageInfo

func (m *LoginQueueMsg) GetPos() int32 {
	if m != nil {
		return m.Pos
	}
	return 0
}

func (m *LoginQueueMsg) GetNum() int32 {
	if m != nil {
		return m.Num
	}
	return 0
}

// 心跳
type HeartBeatReq struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatReq) Reset()         { *m = HeartBeatReq{} }
func (m *HeartBeatReq) String() string { return proto.CompactTextString(m) }
func (*HeartBeatReq) ProtoMessage()    {}
func (*HeartBeatReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{6}
}

func (m *HeartBeatReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatReq.Unmarshal(m, b)
}
func (m *HeartBeatReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatReq.Marshal(b, m, deterministic)
}
func (m *HeartBeatReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatReq.Merge(m, src)
}
func (m *HeartBeatReq) XXX_Size() int {
	return xxx_messageInfo_HeartBeatReq.Size(m)
}
func (m *HeartBeatReq) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatReq.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatReq proto.InternalMessageInfo

// 心跳回复
type HeartBeatAck struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartBeatAck) Reset()         { *m = HeartBeatAck{} }
func (m *HeartBeatAck) String() string { return proto.CompactTextString(m) }
func (*HeartBeatAck) ProtoMessage()    {}
func (*HeartBeatAck) Descriptor() ([]byte, []int) {
	return fileDescriptor_c06e4cca6c2cc899, []int{7}
}

func (m *HeartBeatAck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartBeatAck.Unmarshal(m, b)
}
func (m *HeartBeatAck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartBeatAck.Marshal(b, m, deterministic)
}
func (m *HeartBeatAck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartBeatAck.Merge(m, src)
}
func (m *HeartBeatAck) XXX_Size() int {
	return xxx_messageInfo_HeartBeatAck.Size(m)
}
func (m *HeartBeatAck) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartBeatAck.DiscardUnknown(m)
}

var xxx_messageInfo_HeartBeatAck proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("msg.MSGID", MSGID_name, MSGID_value)
	proto.RegisterType((*LoginReq)(nil), "msg.LoginReq")
	proto.RegisterType((*LoginAck)(nil), "msg.LoginAck")
	proto.RegisterType((*LoginMsg)(nil), "msg.LoginMsg")
	proto.RegisterType((*LoginQueueReq)(nil), "msg.LoginQueueReq")
	proto.RegisterType((*LoginQueueAck)(nil), "msg.LoginQueueAck")
	proto.RegisterType((*LoginQueueMsg)(nil), "msg.LoginQueueMsg")
	proto.RegisterType((*HeartBeatReq)(nil), "msg.HeartBeatReq")
	proto.RegisterType((*HeartBeatAck)(nil), "msg.HeartBeatAck")
}

func init() {
	proto.RegisterFile("msg.proto", fileDescriptor_c06e4cca6c2cc899)
}

var fileDescriptor_c06e4cca6c2cc899 = []byte{
	// 259 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0xc1, 0x4b, 0xc3, 0x30,
	0x18, 0xc5, 0xed, 0x62, 0xc7, 0xf2, 0xe1, 0xe6, 0xc7, 0x87, 0x48, 0x0f, 0x1e, 0x4a, 0x4f, 0xe2,
	0xc1, 0xcb, 0x2e, 0x1e, 0xad, 0x88, 0x53, 0x5c, 0x0f, 0xc6, 0x9b, 0xb7, 0x59, 0x43, 0x18, 0xb5,
	0xcd, 0xda, 0x24, 0xf7, 0xfe, 0xe9, 0x92, 0xac, 0x6c, 0x05, 0x41, 0xf0, 0xf6, 0xde, 0x4b, 0x7e,
	0x8f, 0x07, 0x1f, 0xf0, 0xda, 0xa8, 0xdb, 0x5d, 0xa7, 0xad, 0x26, 0x56, 0x1b, 0x95, 0x5d, 0xc1,
	0x6c, 0xad, 0xd5, 0xb6, 0x11, 0xb2, 0x25, 0x04, 0xe6, 0xb6, 0x5f, 0x49, 0x94, 0x46, 0xd7, 0x5c,
	0x78, 0x99, 0xdd, 0x0d, 0xaf, 0x79, 0x59, 0xd1, 0x25, 0x4c, 0x3b, 0x69, 0xdc, 0xb7, 0x0d, 0x1f,
	0x62, 0x31, 0x38, 0xba, 0x80, 0xd8, 0xea, 0x4a, 0x36, 0xc9, 0x24, 0x70, 0x7b, 0x73, 0x20, 0x0b,
	0xa3, 0xfe, 0x49, 0x9e, 0xc3, 0x3c, 0x90, 0x6f, 0x4e, 0x3a, 0x29, 0x64, 0x9b, 0xbd, 0x8e, 0x83,
	0xbf, 0x96, 0x20, 0xb0, 0x9d, 0x36, 0xa1, 0x2d, 0x16, 0x5e, 0xfa, 0xa4, 0x71, 0x75, 0xc2, 0xf6,
	0x49, 0xe3, 0xea, 0x6c, 0x39, 0x2e, 0xf3, 0xe3, 0x06, 0x28, 0xfa, 0x05, 0x4d, 0x8e, 0xd0, 0x02,
	0xce, 0x9e, 0xe5, 0xa6, 0xb3, 0x0f, 0x72, 0x63, 0xfd, 0xa2, 0xb1, 0xcf, 0xcb, 0xea, 0xe6, 0x03,
	0xe2, 0xe2, 0x7d, 0xf5, 0xf2, 0x48, 0x33, 0x38, 0x7d, 0x5a, 0xe7, 0x2b, 0x3c, 0xa1, 0x19, 0x30,
	0x21, 0x5b, 0xec, 0x53, 0xaf, 0xf2, 0xb2, 0xc2, 0xfe, 0x9e, 0x38, 0xb0, 0xc2, 0x28, 0xec, 0xfb,
	0x88, 0x38, 0xc4, 0x61, 0x06, 0x22, 0x2d, 0x00, 0x8e, 0x8b, 0x90, 0x68, 0x0e, 0xfc, 0x50, 0x8e,
	0xe9, 0xe7, 0x34, 0x1c, 0x6b, 0xf9, 0x13, 0x00, 0x00, 0xff, 0xff, 0x4e, 0x7d, 0x6a, 0x00, 0xb9,
	0x01, 0x00, 0x00,
}

package utils

import (
	"QueueService/msg"
	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true, TimeFormat: "2006-01-02 15:04:05.000"})
}

func Test_MsgReflect(t *testing.T) {
	log.Info().Msg("Test_MsgReflect")
	var req msg.LoginReq
	req.Uid = "123"
	req.Token = "456"
	log.Info().Msg(req.String())
	data, _ := proto.Marshal(&req)

	req2, _ := TheMsgMgr.Unmarshal(123, data)

	log.Info().Msg(req2.(*msg.LoginReq).String())
}

func Test_Token(t *testing.T) {
	log.Info().Msg("Test_Token")
	log.Info().Msg(GetToken())
}

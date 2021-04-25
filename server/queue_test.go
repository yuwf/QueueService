package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
	"time"
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true, TimeFormat: "2006-01-02 15:04:05.000"})
}

func Test_queue(t *testing.T) {
	theQueue.Init()

	theQueue.AddClient("123", nil)
	theQueue.AddClient("456", nil)
	theQueue.AddClient("789", nil)

	time.Sleep(time.Second * 2)

	{
		pos, len := theQueue.GetClientPos("456")
		log.Info().Int("pos", pos).Int("len", len).Msg("123 pos")
	}

	theQueue.RemoveClient("456")
	theQueue.AddClient("111", nil)
	theQueue.AddClient("222", nil)
	theQueue.RemoveClient("111")
	time.Sleep(time.Second * 2)

	{
		pos, len := theQueue.GetClientPos("111")
		log.Info().Int("pos", pos).Int("len", len).Msg("111 pos")
	}
}

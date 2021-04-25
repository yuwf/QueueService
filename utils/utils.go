package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"math/rand"
	"runtime"
	"time"
)

var tokenSeed = []byte{
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9',
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

var tokenSeedSize = len(tokenSeed)

//生成Token字串
func GetToken() string {
	rand.Seed(time.Now().UnixNano())
	token := []byte{}
	for i := 0; i < 32; i++ {
		index := rand.Intn(tokenSeedSize)
		token = append(token, tokenSeed[index])
	}
	return string(token)
}

//生成Token字串
func GetUserId() string {
	name := []byte{}
	for i := 0; i < 5; i++ {
		index := rand.Intn(tokenSeedSize)
		name = append(name, tokenSeed[index])
	}
	return "User_" + string(name)
}

func HandlePanic() {
	if r := recover(); r != nil {
		buf := make([]byte, 2048)
		l := runtime.Stack(buf, false)
		err := fmt.Errorf("%v: %s", r, buf[:l])
		log.Error().Err(err).Msg("Panic")
	}
}

package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"strings"

	"github.com/LCBHSStudent/xfw-core/src/poet"
	"github.com/LCBHSStudent/xfw-core/src/random-gck"
)

type simpleFunc 		func () string
type idGroupFunc 		func (int64, int64)
type idGroupMsgFunc 	func (int64, int64, string)
type idMsgFunc 			func (int64, string)
type groupMsgFunc 		func (int64, string)

var simpleFuncRouter map[string] simpleFunc

func init() {
	simpleFuncRouter = make(map[string] simpleFunc)
	simpleFuncRouter["xfw"] = poet.GetPoetry
	simpleFuncRouter["XFW"] = poet.GetPoetry
	simpleFuncRouter["小飞舞"] = poet.GetPoetry
}

func routeByPrefix(msg string) (groupMsgFunc, int, string) {
	if strings.HasPrefix(msg, "称呼+") {
		return randomGck.SaveAddress, 7, "已添加称呼:" + msg[7:]
	} else if strings.HasPrefix(msg, "形容+") {
		return randomGck.SaveDescription, 7, "已添加形容:" + msg[7:]
	}
	return nil, -1, ""
}

func randomTrigger() simpleFunc {
	bProb, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}
	prob := int(bProb.Uint64())

	if prob <= 4 {
		return randomGck.GenerateSpeech
	} else {
		return nil
	}
}

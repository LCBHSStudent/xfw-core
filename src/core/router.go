package main

import (
	"crypto/rand"
	"log"
	"math/big"
	"regexp"
	"strings"

	"github.com/LCBHSStudent/xfw-core/src/poet"
	randomGck "github.com/LCBHSStudent/xfw-core/src/random-gck"
	"github.com/LCBHSStudent/xfw-core/util"
)

type simpleFunc func() string
type idGroupFunc func(int64, int64)
type idGroupMsgFunc func(int64, int64, string)
type idMsgFunc func(int64, string)
type groupMsgFunc func(int64, string)

var simpleFuncRouter map[string]simpleFunc

var collectRandomly = true

var (
	考研Regexp *regexp.Regexp
	地域Regexp *regexp.Regexp
	学历Regexp *regexp.Regexp
)

func InitRegexps() {
	考研Regexp = regexp.MustCompile(`(.*?)考研(.*?)`)
	地域Regexp = regexp.MustCompile(`(.*?)(沙东|山东)(.*?)`)
	学历Regexp = regexp.MustCompile(`(.*?)(北邮|学历)(.*?)`)
}

func init() {
	InitRegexps()

	simpleFuncRouter = make(map[string]simpleFunc)
	simpleFuncRouter["xfw"] = poet.GetPoetry
	simpleFuncRouter["XFW"] = poet.GetPoetry
	simpleFuncRouter["小飞舞"] = poet.GetPoetry
	simpleFuncRouter["龚诗"] = poet.GenerateGongPoem
	simpleFuncRouter["小万邦"] = poet.GenerateGongPoem
	simpleFuncRouter["Collect Randomly=ON"] = func() string {
		collectRandomly = true
		return "自动收集已启动"
	}
	simpleFuncRouter["Collect Randomly=OFF"] = func() string {
		collectRandomly = false
		return "自动收集已关闭"
	}

}

func routeByPrefix(msg string) (groupMsgFunc, int, string) {
	if strings.HasPrefix(msg, "称呼+") {
		return randomGck.SaveAddress, 7, "已添加称呼:" + msg[7:]
	} else if strings.HasPrefix(msg, "形容+") {
		return randomGck.SaveDescription, 7, "已添加形容:" + msg[7:]
	}
	return nil, -1, ""
}

func routeBy学历地域工作出身(msg string) (int, string) {
	var ret string
	matchs := 考研Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		ret += "考研的"
	}

	matchs = 地域Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		ret += "沙东人"
	}

	matchs = 学历Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		ret += "明明就是个臭带专的，"
	}

	if ret == "" {
		return -1, ""
	}

	return 0, ret + randomGck.GenerateDescription()
}

func randomTrigger(targetId int64, msg string) (ret simpleFunc) {
	bProb, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}
	prob := int(bProb.Uint64())

	if collectRandomly && (prob <= 2 || len(msg) > 4) {
		randomGck.SaveDescription(targetId, msg)
	}

	if prob < 2 {
		if _, ok := util.GetObjectByKey("group-enable-send-randomgck").(map[int64]bool)[targetId]; ok {
			ret = randomGck.GenerateSpeech
		}
	}

	return
}

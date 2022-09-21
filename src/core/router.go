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
	地域Regexp *regexp.Regexp
	学历Regexp *regexp.Regexp
	工作Regexp *regexp.Regexp
)

func InitRegexps() {
	地域Regexp = regexp.MustCompile(`(.*?)(沙东|北京|百京|成都|长沙|通化|宜春|江西|东百|东北|河南|荷兰|上海|郴州|青岛|新疆|武汉|湖北|天津|山东|白银|贵阳|香港|HK|加拿大|海南|深圳|青浦)(.*?)`)
	学历Regexp = regexp.MustCompile(`(.*?)(北邮|考研|保研|本科|留学|直博|phd)(.*?)`)
	工作Regexp = regexp.MustCompile(`(.*?)(智加|字节|图森|泰康|滴滴|tp|银行|很行)(.*?)`)
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
	var 学历Prefix string
	var 地域Prefix string
	var 工作Prefix string

	cqcodeExp := regexp.MustCompile(`\[CQ:[\s\S]*?\]`)
	msg = cqcodeExp.ReplaceAllString(msg, "")

	matchs := 学历Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		prob := util.GetRandNum(10)
		学历Prefix = "明明就是个臭带专的" + matchs[2] + "逼"
		if prob <= 2 {
			学历Prefix = "某些特别优秀的" + matchs[2] + "爹"
		}
	}

	matchs = 地域Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		prob := util.GetRandNum(10)

		suffix := "人"
		if prob <= 2 {
			suffix = "逼"
		} else if prob >= 9 {
			suffix = "爹"
		}

		地域Prefix = matchs[2] + suffix
	}

	matchs = 工作Regexp.FindStringSubmatch(msg)
	if len(matchs) > 0 {
		prob := util.GetRandNum(100)

		suffix := "爹"
		if prob <= 20 {
			suffix = "逼"
		}

		if prob >= 95 {
			suffix += "失业后"
		}

		工作Prefix = matchs[2] + suffix
	}

	ret := 学历Prefix + 地域Prefix + 工作Prefix
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

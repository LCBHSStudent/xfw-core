package main

import (
	"log"
	"strings"

	qqbotapi "github.com/catsworld/qq-bot-api"
	util "github.com/LCBHSStudent/xfw-core/util"
	_ "github.com/LCBHSStudent/xfw-core/src/database"
)

const CQHttpConnKey = "cqhttp-ws-connect"

func main() {
	var bot *qqbotapi.BotAPI
	var err error
	
	cqhttpConf := util.GetObjectByKey(CQHttpConnKey).(map[interface{}]interface{})

	bot, err = qqbotapi.NewBotAPI("",
		strings.Join([]string{
			"ws://",
			cqhttpConf["ipv4"].(string),
			":",
			cqhttpConf["port"].(string),
		}, ""), cqhttpConf["secret"].(string),
	)

	if err != nil {
		log.Fatal(err)
	}
	bot.Debug = false

	conf := qqbotapi.NewUpdate(0)
	conf.PreloadUserInfo = true
	updates, err := bot.GetUpdatesChan(conf)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if handle, ok := simpleFuncRouter[update.Message.Text]; ok {
			targetId := update.GroupID
			if update.MessageType == "private" {
				targetId = update.UserID
			}
			
			bot.NewMessage(targetId, update.MessageType).Text(handle()).Send()
		}
	}
}

package main

import (
	"log"
	"strings"

	qqbotapi "github.com/catsworld/qq-bot-api"
)

const CqHttpHostWS = "127.0.0.1:6700"

func main() {
	var bot *qqbotapi.BotAPI

	var err error
	bot, err = qqbotapi.NewBotAPI("", strings.Join([]string{"ws://", CqHttpHostWS}, ""), "CQHTTP_SECRET")

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
	}
}

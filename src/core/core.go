package main

import (
	"log"
	"strconv"
	"strings"

	randomGck "github.com/LCBHSStudent/xfw-core/src/random-gck"
	util "github.com/LCBHSStudent/xfw-core/util"
	qqbotapi "github.com/catsworld/qq-bot-api"
	cqcode "github.com/catsworld/qq-bot-api/cqcode"
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

		isGroupMsg := update.MessageType == "group"
		isPrivateMsg := update.MessageType == "private"
		fromIdStr := strconv.FormatInt(update.Message.From.ID, 10)

		targetId := update.GroupID
		if isPrivateMsg {
			targetId = update.UserID
		}

		if bot.IsMessageToMe(*update.Message) {
			message := make(cqcode.Message, 0)
			message.Append(&cqcode.At{QQ: fromIdStr})
			message.Append(&cqcode.Text{Text: "\n"})

			if !checkUserBlackList(update.Message.From.ID) {
				parseRichMessage(randomGck.GenerateSpeech(), &message)
			} else {
				message.Append(&cqcode.Text{Text: "给我GCK啊！！！！！"})
			}
			bot.SendMessage(update.GroupID, "group", message)
			continue
		}

		// random triggered function
		if handle := randomTrigger(update.GroupID, update.Message.Text); handle != nil {
			if isGroupMsg {
				message := make(cqcode.Message, 0)
				message.Append(&cqcode.At{QQ: fromIdStr})
				message.Append(&cqcode.Text{Text: "\n"})

				parseRichMessage(handle(), &message)
				bot.SendMessage(update.GroupID, "group", message)
				continue
			}
		}

		if handle, ok := simpleFuncRouter[update.Message.Text]; ok {
			if checkUserBlackList(update.Message.From.ID) {
				bot.NewMessage(targetId, update.MessageType).Text("GCK！！！！！").Send()
			} else {
				go func() {
					bot.NewMessage(targetId, update.MessageType).Text(handle()).Send()
				}()
			}
		} else if handle, ok, msg := routeByPrefix(update.Message.Text); ok >= 0 {
			if isGroupMsg {
				message := make(cqcode.Message, 0)
				message.Append(&cqcode.At{QQ: fromIdStr})
				message.Append(&cqcode.Text{Text: "\n"})
				if !checkUserBlackList(update.Message.From.ID) {
					handle(update.GroupID, update.Message.Text[ok:])
					parseRichMessage(msg, &message)
				} else {
					message.Append(&cqcode.Text{Text: "给我GCK啊！！！！！"})
				}
				bot.SendMessage(update.GroupID, "group", message)
			}
		} else if ok, msg := routeBy学历地域工作出身(update.Message.Text); ok >= 0 {
			if isGroupMsg {
				message := make(cqcode.Message, 0)

				if _, ok := util.GetObjectByKey("group-enable-send-randomgck").(map[int64]bool)[targetId]; ok {
					if !checkUserBlackList(update.Message.From.ID) {
						parseRichMessage(msg, &message)
					} else {
						message.Append(&cqcode.Text{Text: "给我GCK啊！！！！！"})
					}

					bot.SendMessage(update.GroupID, "group", message)
				}
			}

		} else {
			if len(msg) != 0 && isGroupMsg {
				bot.NewMessage(update.GroupID, "group").
					At(fromIdStr).
					NewLine().
					Text(msg).Send()
			}
		}

	}
}

var ignoredType = map[string]bool{
	"at":    true,
	"reply": true,
}

func parseRichMessage(raw string, message *cqcode.Message) {
	richMessage, err := cqcode.ParseMessage(raw)
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range richMessage {
		if _, ok := ignoredType[v.FunctionName()]; !ok {
			(*message).Append(v)
		}
	}
}

func checkUserBlackList(userID int64) bool {
	if _, ok := util.GetObjectByKey("user-black-list").(map[int64]bool)[userID]; ok {
		return true
	} else {
		return false
	}
}

package chatgpt

import (
	"strings"

	GPT "github.com/golang-infrastructure/go-ChatGPT"
)

var chat *GPT.ChatGPT = nil

func SetToken(groupId int64, token string) string {
	token = strings.TrimSpace(token)
	chat = GPT.NewChatGPT(token)

	return "Token设置完成"
}

func Communicate(groupId int64, message string) string {
	message = strings.TrimSpace(message)
	if chat == nil {
		return "Chat-GPT模块未初始化，请打开chatGPT页面(https://chat.openai.com/chat)后\n" +
			"在F12控制台中运行以下javascript获取token:\n" +
			"JSON.parse(document.getElementById(\"__NEXT_DATA__\").text).props.pageProps.accessToken"
	}

	talk, err := chat.Talk(message)
	if err != nil {
		return err.Error()
	} else {
		return strings.Join(talk.Message.Content.Parts, "\n")
	}
}

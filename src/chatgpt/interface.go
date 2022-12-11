package chatgpt

import (
	"strings"
)

var chat *ChatGPT = nil

func SetToken(groupId int64, token string) string {
	sessionToken := strings.TrimSpace(token)
	chat = Init(sessionToken)

	return "初始化完成"
}

func Communicate(groupId int64, message string) string {
	message = strings.TrimSpace(message)
	if chat == nil {
		return "请先通过'/set-token'命令设置session-token以初始化该模块"
	}

	ch, err := chat.SendMessage(message, groupId)
	if err != nil {
		return err.Error()
	} else {
		result := ""
	loop:
		for {
			select {
			case new, ok := <-ch:
				if ok {
					result = new.Message
				} else {
					break loop
				}
			}
		}

		return result
	}
}

func ResetConversation(groupId int64, message string) string {
	chat.ResetConversation(groupId)

	return "重置对话串成功"
}

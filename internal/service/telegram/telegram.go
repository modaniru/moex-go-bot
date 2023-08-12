package telegram

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageSender struct{
	bot *tgbotapi.BotAPI
}

func NewMessageSender(bot *tgbotapi.BotAPI) *MessageSender{
	return &MessageSender{bot: bot}
}

func (m *MessageSender) SendMessage(message string, id int) error{
	mes := tgbotapi.NewMessage(int64(id), message)
	_, err := m.bot.Send(mes)
	if err != nil{
		return fmt.Errorf("send message error: %w", err)
	}
	return nil
}
package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/moex-telegram-bot/internal/handler"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	handler *handler.Handler
}

func NewBot(bot *tgbotapi.BotAPI, handler *handler.Handler) *Bot {
	return &Bot{
		bot:     bot,
		handler: handler,
	}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.bot.GetUpdatesChan(u)

	for update := range updates {
		b.handler.HandleAction(update)
	}
}

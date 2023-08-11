package handler

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/moex-telegram-bot/internal/service"
)

type Handler struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
}

func NewHandler(bot *tgbotapi.BotAPI, service *service.Service) *Handler {
	return &Handler{bot: bot, service: service}
}

func (h *Handler) HandleAction(update tgbotapi.Update) {
	message := update.Message
	ctx := context.WithValue(context.Background(), "id", message.From.ID)
	ctx = context.WithValue(ctx, "message", message.Text)

	switch message.Text {
	case "/start":
		h.Start(message)
	case "/help":
		h.Help(message)
	case "/register":
		h.Register(message)
	}
}

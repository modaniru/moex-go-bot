package handler

import (
	"context"
	"strings"

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
	command := strings.Split(message.Text, " ")
	switch command[0] {
	case "/start":
		h.Start(message)
	case "/help":
		h.Help(message)
	case "/register":
		h.Register(message)
	case "/add":
		h.AddSecurity(message)
	case "/follow":
		h.isValidUser(message, h.FollowUser)
	case "/unfollow":
		h.isValidUser(message, h.UnfollowUser)
	}
}

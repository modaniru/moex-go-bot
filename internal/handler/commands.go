package handler

import (
	"database/sql"
	"errors"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (h *Handler) Start(m *tgbotapi.Message) {
	message := tgbotapi.NewMessage(
		m.From.ID,
		"Привет! 👋\n\nБлагодаря этому боту ты сможешь вовремя получать уведомления о состоянии разных бумаг.\n\n/help - для получения информации",
	)
	_, err := h.bot.Send(message)
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
		return
	}
	return
}

func (h *Handler) Help(m *tgbotapi.Message) {
	message := tgbotapi.NewMessage(
		m.From.ID,
		"/register - зарегистрироваться для пользования.\n/unRegister - удалить аккаунт.",
	)
	_, err := h.bot.Send(message)
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
		return
	}
	return
}

func (h *Handler) Register(m *tgbotapi.Message) {
	id := m.From.ID
	_, err := h.service.GetUserById(int(id))
	responseMessage := ""

	defer h.handleResponse(&responseMessage, err, id)

	if errors.Is(err, sql.ErrNoRows) {
		err = h.service.Register(int(id))
		if err == nil {
			responseMessage = "Ты успешно зарегистрировался! 🎉"
		}
	} else if err == nil {
		responseMessage = "Ты уже зарегистрирован. 🤨"
	}
}

func (h *Handler) handleResponse(text *string, err error, id int64){
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
		*text = "Ошибка. 😔"
	}
	message := tgbotapi.NewMessage(
		id,
		*text,
	)
	_, err = h.bot.Send(message)
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
	}
}

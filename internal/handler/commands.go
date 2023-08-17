package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/modaniru/moex-telegram-bot/internal/entity"
	"github.com/modaniru/moex-telegram-bot/internal/service/services"
)

func (h *Handler) Start(m *tgbotapi.Message) {
	text := "Привет! 👋\n\nБлагодаря этому боту ты сможешь вовремя получать уведомления о состоянии разных бумаг.\n\n/help - для получения информации"
	h.handleResponse(text, nil, m.From.ID)
}

func (h *Handler) Help(m *tgbotapi.Message) {
	text := "/register - зарегистрироваться для пользования.\n/unRegister - удалить аккаунт."
	h.handleResponse(text, nil, m.From.ID)
}

func (h *Handler) Register(m *tgbotapi.Message) {
	id := m.From.ID
	_, err := h.service.GetUserById(int(id))
	responseMessage := ""

	defer func() { h.handleResponse(responseMessage, err, id) }()

	if errors.Is(err, sql.ErrNoRows) {
		err = h.service.Register(int(id))
		if err == nil {
			responseMessage = "Ты успешно зарегистрировался! 🎉"
		}
	} else if err == nil {
		responseMessage = "Ты уже зарегистрирован. 🤨"
	}
}

func (h *Handler) AddSecurity(m *tgbotapi.Message) {
	// add engine market boardGroup security baseDate coefficient
	id := m.From.ID
	message := m.Text
	args := strings.Split(message, " ")
	if len(args) < 7 {
		h.handleResponse("Нехватает аргументов! 🚫", nil, id)
		return
	} else if len(args) > 7 {
		h.handleResponse("Cлишком много аргументов! 🚫", nil, id)
		return
	}
	bgId, err := strconv.Atoi(args[3])
	if err != nil {
		h.handleResponse("boardgroup должен быть целочисленным! 🔟", nil, id)
		return
	}
	coefficent, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		h.handleResponse("coefficient должен быть целым/дробным числом! 🔟", nil, id)
		return
	}
	// TODO return security and median struct
	resp, err := h.service.TrackService.SaveTrack(&entity.SaveTrack{
		UserId:       int(id),
		Engine:       args[1],
		Market:       args[2],
		BoardGroupId: bgId,
		Security:     args[4],
		Date:         args[5],
		Interval:     10,
		Coefficient:  coefficent,
	})
	if err != nil {
		if errors.Is(err, services.ErrBadDay) {
			h.handleResponse("Вы выбрали день, в который торги не осуществлялись! 🕟", nil, id)
			return
		}
		h.handleResponse("", err, id)
		return
	}
	respMessage := fmt.Sprintf("Ты сохранил бумагу %s, если она превысит volume %d, мы тебя уведомим! 👀\n\nМаксимальный volume(10 min): %d\nМинимальный volume(10 min): %d\nОтслеживаемая дата: %s", resp.Security, resp.Median, resp.MaxVolume, resp.MinVolume, resp.Date)
	h.handleResponse(respMessage, nil, id)
}

func (h *Handler) FollowUser(m *tgbotapi.Message) {
	id := m.From.ID
	u, err := h.service.GetUserById(int(id))
	message := ""

	defer func() {
		h.handleResponse(message, err, id)
	}()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message = "Пользователь не найден! 😢"
			err = nil
		}
		return
	}
	if u.Followed == true {
		message = "Ты уже отслеживаешь! 🙃"
		return
	}
	err = h.service.FollowUser(int(id))

	message = "Ты успешно начал отслеживать уведомления! 😎"
}

func (h *Handler) UnfollowUser(m *tgbotapi.Message) {
	id := m.From.ID
	u, err := h.service.GetUserById(int(id))
	message := ""

	defer func() {
		h.handleResponse(message, err, id)
	}()

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			message = "Пользователь не найден! 😢"
			err = nil
		}
		return
	}
	if u.Followed == false {
		message = "Ты уже не отслеживаешь! 🙃"
		return
	}
	err = h.service.UnfollowUser(int(id))

	message = "Ты отключил уведомления от бота! 😪"
}

// middleware that validate user
func (h *Handler) isValidUser(m *tgbotapi.Message, next func(*tgbotapi.Message)) {
	id := m.From.ID
	_, err := h.service.GetUserById(int(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.handleResponse("Ты не зарегистрирован! 😲", nil, id)
			return
		}
		h.handleResponse("", err, id)
	} else {
		next(m)
	}
}

// Send message to user based in error
func (h *Handler) handleResponse(text string, err error, id int64) {
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
		text = "Ошибка. 😔"
	}
	message := tgbotapi.NewMessage(
		id,
		text,
	)
	_, err = h.bot.Send(message)
	if err != nil {
		slog.Error("send message error", slog.String("error", err.Error()))
	}
}

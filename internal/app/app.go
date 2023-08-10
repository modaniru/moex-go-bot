package app

import (
	"log/slog"
	"os"

	"github.com/modaniru/moex-telegram-bot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// https://iss.moex.com/iss/engines/stock/markets/shares/boardgroups/57/securities/SBER/candles.jsonp?interval=1&from=2023-08-10
func App() {
	cfg := config.LoadConfig()
	ConfigureLogger(cfg.Level)
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		slog.Error("create bot api error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message.Text == "kek" {
			m := tgbotapi.NewMessage(update.Message.From.ID, "lol")
			bot.Send(m)
		}
	}
}

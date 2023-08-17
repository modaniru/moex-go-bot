package app

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/modaniru/moex-telegram-bot/config"
	"github.com/modaniru/moex-telegram-bot/internal/bot"
	"github.com/modaniru/moex-telegram-bot/internal/clients"
	"github.com/modaniru/moex-telegram-bot/internal/handler"
	"github.com/modaniru/moex-telegram-bot/internal/service"
	"github.com/modaniru/moex-telegram-bot/internal/service/notifier"
	"github.com/modaniru/moex-telegram-bot/internal/service/telegram"
	"github.com/modaniru/moex-telegram-bot/internal/storage"
	"github.com/modaniru/moex-telegram-bot/internal/storage/gen"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// https://iss.moex.com/iss/engines/stock/markets/shares/boardgroups/57/securities/SBER/candles.jsonp?interval=1&from=2023-08-10
func App() {
	cfg := config.LoadConfig()
	ConfigureLogger(cfg.Level)
	slog.Info("create bot api...")
	botApi, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		slog.Error("create bot api error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Authorized on account " + botApi.Self.UserName)
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:1111/postgres?sslmode=disable")
	if err != nil {
		slog.Error("open db error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = db.Ping()
	if err != nil {
		slog.Error("ping db error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("init gen.Queries...")
	q := gen.New(db)

	slog.Info("init http.Client...")
	client := http.Client{}

	slog.Info("init clients.MoexClient...")
	moex := clients.NewMoexClient(&client)

	slog.Info("init storage...")
	storage := storage.NewStorage(db, q)

	slog.Info("init service...")
	s := service.NewService(moex, storage)

	slog.Info("init handler...")
	handler := handler.NewHandler(botApi, s)

	slog.Info("init message sender...")
	sender := telegram.NewMessageSender(botApi)

	slog.Info("init notifier...")
	notifier := notifier.NewNotifier(s.TrackService, sender)

	slog.Info("init bot server...")
	botServer := bot.NewBot(botApi, handler)

	slog.Info("start notifier...")
	notifier.StartNotifier()

	slog.Info("start bot server...")
	botServer.Start()
	slog.Error("bot was not up")
}

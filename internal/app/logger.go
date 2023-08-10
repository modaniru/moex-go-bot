package app

import (
	"log/slog"
	"os"
)

const (
	debug = "debug"
	dev   = "dev"
	prod  = "prod"
)

func ConfigureLogger(level string) {
	var handler slog.Handler
	switch level {
	case debug:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	case dev:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	default:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	slog.SetDefault(slog.New(handler))
	slog.Debug("logger successfully loaded", slog.String("level", level))
}

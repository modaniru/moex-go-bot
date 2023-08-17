package notifier

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/modaniru/moex-telegram-bot/internal/service"
	"github.com/modaniru/moex-telegram-bot/internal/service/services"
	"github.com/modaniru/moex-telegram-bot/internal/service/telegram"
)

type Notifier struct {
	candles service.TrackService
	sender  *telegram.MessageSender
}

// TODO sender service
func NewNotifier(candles service.TrackService, message *telegram.MessageSender) *Notifier {
	return &Notifier{
		candles: candles,
		sender:  message,
	}
}

const (
	up   = "↗️"
	down = "↘️"
)

func (n *Notifier) StartNotifier() {
	for time.Now().Second() > 5 {
		slog.Debug("we are trying to get closer to the exact data", slog.Int("delta from start minute in second", time.Now().Second()))
		time.Sleep(3 * time.Second)
	}
	slog.Debug("success", slog.Int("delta from start minute in second", time.Now().Second()))
	go func() {
		for true {
			time.Sleep(time.Minute)
			rows, err := n.candles.GetAllMustNotifiedTracks()
			if err != nil {
				slog.Error("get all must notifiend candles error", slog.String("err", err.Error()))
				continue
			}
			for _, row := range rows {
				candle, err := n.candles.GetCandle(row)
				if err != nil {
					if errors.Is(err, services.ErrBadDay) {
						continue
					}
					slog.Error("get candle request error", slog.String("err", err.Error()))
					continue
				}
				if candle.Volume >= int(row.TrackedVolume) {
					s := up
					if candle.Open > candle.Close {
						s = down
					}
					n.sender.SendMessage(fmt.Sprintf("volume по бумаге %s превышает %d, сейчас %d. status %s", row.Security, row.TrackedVolume, candle.Volume, s), int(row.ID))
				}

			}

		}
	}()
}

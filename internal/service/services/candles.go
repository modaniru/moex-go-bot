package services

import (
	"fmt"
	"strconv"

	"github.com/modaniru/moex-telegram-bot/internal/clients"
)

type CandlesService struct{
	moexClient *clients.MoexClient
}

func (c *CandlesService)SaveTrack(params *clients.CandleRequest) error{
	params.Interval = 10 //default interval
	candles, err := c.moexClient.Candles(params)
	if err != nil{
		return fmt.Errorf("get candles error %w", err)
	}
	sum := 0
	volumeIndex := 0
	for _, column := range candles.Candles.Columns{
		if column == "volume"{
			break
		}
		volumeIndex++
	}
	for _, arr := range candles.Candles.Data{
		v, err := strconv.Atoi(arr[volumeIndex])
		if err != nil{
			return fmt.Errorf("convert volume error %w", err)
		}
		sum += v
	}
	// median := sum / len(candles.Candles.Data)
	return nil
}
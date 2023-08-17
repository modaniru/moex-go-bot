package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/modaniru/moex-telegram-bot/internal/clients"
	"github.com/modaniru/moex-telegram-bot/internal/entity"
	"github.com/modaniru/moex-telegram-bot/internal/storage/gen"
)

var (
	ErrBadDay = errors.New("bad day, data's len equals zero")
)

type TrackStorage interface {
	UntrackSecurityByUserIdAndId(ctx context.Context, arg gen.UntrackSecurityByUserIdAndIdParams) error
	TrackSecurityByUserIdAndId(ctx context.Context, arg gen.TrackSecurityByUserIdAndIdParams) error
	SaveTrack(ctx context.Context, arg gen.SaveTrackParams) error
	DeleteUser(ctx context.Context, id int32) error
	GetAllMustNotifiedTracks(ctx context.Context) ([]gen.GetAllMustNotifiedTracksRow, error)
	DeleteTrackByUserIdAndId(ctx context.Context, arg gen.DeleteTrackByUserIdAndIdParams) error
	GetUserTracks(ctx context.Context, userID int32) ([]gen.Track, error)
}

type CandlesService struct {
	moexClient *clients.MoexClient
	storage    TrackStorage
}

// TODO избавиться от создания множества структур
func NewCandlesService(moex *clients.MoexClient, storage TrackStorage) *CandlesService {
	return &CandlesService{moexClient: moex, storage: storage}
}

func (c *CandlesService) SaveTrack(params *entity.SaveTrack) (*entity.TrackResponse, error) {
	candles, err := c.moexClient.Candles(&clients.CandleRequest{
		Engine:       params.Engine,
		Market:       params.Market,
		BoardGroupId: params.BoardGroupId,
		Security:     params.Security,
		Date:         params.Date,
		Interval:     params.Interval,
	})
	if err != nil {
		return nil, fmt.Errorf("get candles error %w", err)
	}
	if len(candles.Candles.Data) == 0 {
		return nil, ErrBadDay
	}
	sum := 0
	volumeIndex := 0
	for _, column := range candles.Candles.Columns {
		if column == "volume" {
			break
		}
		volumeIndex++
	}
	min, ok := candles.Candles.Data[0][volumeIndex].(float64)
	if !ok {
		return nil, fmt.Errorf("convert volume error")
	}
	max, ok := candles.Candles.Data[0][volumeIndex].(float64)
	if !ok {
		return nil, fmt.Errorf("convert volume error")
	}
	for _, arr := range candles.Candles.Data {
		v, ok := arr[volumeIndex].(float64)
		if !ok {
			return nil, fmt.Errorf("convert volume error")
		}
		if v > max {
			max = v
		} else if v < min {
			min = v
		}
		sum += int(v)
	}
	median := int(float64(sum/(len(candles.Candles.Data)*params.Interval)) * params.Coefficient)
	err = c.storage.SaveTrack(context.Background(), gen.SaveTrackParams{
		UserID:        int32(params.UserId),
		Engine:        params.Engine,
		Market:        params.Market,
		BoardGroup:    int32(params.BoardGroupId),
		Security:      params.Security,
		Date:          params.Date,
		TrackedVolume: int32(median),
	})
	if err != nil {
		return nil, fmt.Errorf("save track error: %w", err)
	}
	return &entity.TrackResponse{
		Median:    median,
		Security:  params.Security,
		MaxVolume: int(max),
		MinVolume: int(min),
		Date:      params.Date,
	}, nil
}

func (c *CandlesService) GetAllMustNotifiedTracks() ([]gen.GetAllMustNotifiedTracksRow, error) {
	return c.storage.GetAllMustNotifiedTracks(context.Background())
}

type CandleResponse struct {
	Open   int
	Close  int
	Volume int
}

func (c *CandlesService) GetCandle(row gen.GetAllMustNotifiedTracksRow) (*CandleResponse, error) {
	resp, err := c.moexClient.Candles(
		&clients.CandleRequest{
			Engine:       row.Engine,
			Market:       row.Market,
			BoardGroupId: int(row.BoardGroup),
			Security:     row.Security,
			Date:         time.Now().Format("2006-01-02"),
			Interval:     1,
			IsReverse:    true,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("send request error: %w", err)
	}
	if len(resp.Candles.Data) == 0 {
		return nil, ErrBadDay
	}
	// TODO search by columns

	volume, ok := resp.Candles.Data[0][5].(float64)
	if !ok {
		return nil, fmt.Errorf("conv value error")
	}
	open, ok := resp.Candles.Data[0][0].(float64)
	if !ok {
		return nil, fmt.Errorf("conv value error")
	}
	close, ok := resp.Candles.Data[0][1].(float64)
	if !ok {
		return nil, fmt.Errorf("conv value error")
	}

	return &CandleResponse{
		Volume: int(volume),
		Open:   int(open),
		Close:  int(close),
	}, nil
}

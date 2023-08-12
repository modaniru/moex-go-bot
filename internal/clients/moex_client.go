package clients

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	url = "https://iss.moex.com/iss/engines/%s/markets/%s/boardgroups/%d/securities/%s/candles.jsonp?interval=%d&from=%s&till=%s"
)

type MoexClient struct {
	client *http.Client
}

func NewMoexClient(client *http.Client) *MoexClient {
	return &MoexClient{client: client}
}

func (m *MoexClient) Candles(req *CandleRequest) (*CandlesResponse, error) {
	url := fmt.Sprintf(url, req.Engine, req.Market, req.BoardGroupId, req.Security, req.Interval, req.Date, req.Date)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %w", err)
	}
	response, err := m.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("client call request error: %w", err)
	}
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %w", err)
	}
	var result CandlesResponse
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, fmt.Errorf("unmarshal json error: %w", err)
	}
	return &result, nil
}

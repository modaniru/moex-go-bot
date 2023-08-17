package clients

type CandlesResponse struct {
	Candles Candles `json:"candles"`
}

type Candles struct {
	Metadata Metadata        `json:"metadata"`
	Columns  []string        `json:"columns"`
	Data     [][]interface{} `json:"data"`
}

type Metadata struct {
	Open   TypeResponse `json:"open"`
	Close  TypeResponse `json:"close"`
	High   TypeResponse `json:"high"`
	Low    TypeResponse `json:"low"`
	Value  TypeResponse `json:"value"`
	Volume TypeResponse `json:"volume"`
	Begin  TimeResponse `json:"start"`
	End    TimeResponse `json:"end"`
}

type TypeResponse struct {
	Type string `json:"type"`
}

type TimeResponse struct {
	Type    string `json:"type"`
	Bytes   int    `json:"bytes"`
	MaxSize int    `json:"max_size"`
}

type CandleRequest struct {
	Engine       string
	Market       string
	BoardGroupId int
	Security     string
	Date         string // YYYY-MM-DD
	Interval     int    // 1, 10, 60
	IsReverse    bool
}

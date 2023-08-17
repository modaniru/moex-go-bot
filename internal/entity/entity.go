package entity

type SaveTrack struct {
	UserId       int
	Engine       string
	Market       string
	BoardGroupId int
	Security     string
	Date         string // YYYY-MM-DD
	Interval     int    // 1, 10, 60
	Coefficient  float64
}

type TrackResponse struct {
	Security  string
	Median    int
	MinVolume int
	MaxVolume int
	Date      string
}

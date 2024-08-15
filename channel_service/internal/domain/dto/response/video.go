package response

type Video struct {
	VideoID   int       `json:"video_id"`
	ChannelID int       `json:"channel_id"`
	Title     string    `json:"title"`
	Comments  []Comment `json:"comments"`
	IsActive  bool      `json:"is_active"`
}
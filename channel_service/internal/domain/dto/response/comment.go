package response

type Comment struct {
	CommendID      int     `json:"comment_id"`
	CommentContent string  `json:"content"`
	CommentType    string  `json:"comment_type"`
	Price          float64 `json:"price"`
	VideoID        int     `json:"video_id"`
}
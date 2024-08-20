package domain

type Channel struct {
	ChannelID       int     `json:"channel_id"`
	UserID          int     `json:"user_id" binding:"required" validate:"gt=0"`
	ChannelName     string  `json:"channel_name" binding:"required"`
	Membership      []int   `json:"membership"`
	MembershipPrice float64 `json:"membership_price" binding:"required" validate:"gt=0"`
}

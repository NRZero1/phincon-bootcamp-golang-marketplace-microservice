package response

type Channel struct {
	ChannelID       int     `json:"channel_id"`
	UserID          int     `json:"user_id"`
	ChannelName     string  `json:"channel_name"`
	Membership      []int   `json:"membership"`
	MembershipPrice float64 `json:"membership_price"`
	Videos          []Video `json:"videos"`
}

package request

type MembershipRequest struct {
	ChannelID        int  `json:"channel_id"`
	UserIDDest       int  `json:"user_id_dest"`
	IsPaymentPending bool `json:"is_payment_pending"`
}

package request

type DonationRequest struct {
	ProviderID       int  `json:"provider_id"`
	ChannelID        int  `json:"channel_id"`
	UserIDDest       int  `json:"user_id_dest"`
	IsPaymentPending bool `json:"is_payment_pending"`
}

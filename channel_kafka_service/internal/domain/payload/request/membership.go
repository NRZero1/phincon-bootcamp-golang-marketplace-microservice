package request

type MembershipRequest struct {
	ChannelID        int
	UserIDDest       int
	IsPaymentPending bool
}

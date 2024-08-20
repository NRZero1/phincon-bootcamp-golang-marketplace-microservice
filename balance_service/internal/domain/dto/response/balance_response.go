package response

type BalanceResponse struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

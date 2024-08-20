package domain

type Balance struct {
	UserID  int     `json:"user_id"`
	Balance float64 `json:"balance"`
}

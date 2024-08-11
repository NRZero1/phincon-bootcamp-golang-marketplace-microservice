package response

type UserResponse struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	Balance   float64 `json:"balance"`
	PackageID int     `json:"package_id"`
}

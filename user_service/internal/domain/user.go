package domain

type User struct {
	UserID    int     `json:"user_id"`
	Username  string  `json:"username"`
	Password  string  `json:"password"`
	Role      string  `json:"role"`
	Balance   float64 `json:"balance"`
	PackageID int     `json:"package_id"`
}

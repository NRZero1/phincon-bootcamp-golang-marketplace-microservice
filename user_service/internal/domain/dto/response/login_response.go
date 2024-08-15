package response

type LoginResponse struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

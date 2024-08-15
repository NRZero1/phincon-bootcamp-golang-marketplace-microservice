package request

type Register struct {
	Username        string  `json:"username" binding:"required"`
	Password        string  `json:"password" binding:"required"`
	ConfirmPassword string  `json:"confirm_password" binding:"required" validate:"eqfield=password"`
	Balance         float64 `json:"balance"`
}
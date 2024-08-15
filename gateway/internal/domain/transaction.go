package domain

type Transaction struct {
	ID            int    `json:"id,omitempty" binding:"required"`
	TransactionID string `json:"transaction_id,omitempty" binding:"required"`
	OrderType     string `json:"order_type" binding:"required"`
	Status        string `json:"status" binding:"required"`
	UserID        int    `json:"user_id" binding:"required"`
}
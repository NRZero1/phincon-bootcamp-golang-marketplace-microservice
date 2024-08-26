package domain

type OrderRequest struct {
	OrderType     string      `json:"order_type" binding:"required"`
	TransactionID string      `json:"transaction_id"`
	UserID        int         `json:"user_id" binding:"required"`
	Payload       interface{} `json:"payload" binding:"required"`
}

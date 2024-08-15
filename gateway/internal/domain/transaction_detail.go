package domain

type TransactionDetail struct {
	ID            int         `json:"id" binding:"required"`
	TransactionID string      `json:"transaction_id" binding:"required"`
	OrderType     string      `json:"order_type" binding:"required"`
	Topic         string      `json:"topic" binding:"required"`
	Step          string      `json:"step" binding:"required"`
	Service       string      `json:"service" binding:"required"`
	Status        string      `json:"status" binding:"required"`
	Message       string      `json:"message" binding:"required"`
	Payload       interface{} `json:"payload" binding:"required"`
}
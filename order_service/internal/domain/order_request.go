package domain

type OrderRequest struct {
	OrderType     string      `json:"orderType" binding:"required"`
	TransactionID string      `json:"transactionId"`
	UserID        int         `json:"userId" binding:"required"`
	Payload       interface{} `json:"payload" binding:"required"`
}

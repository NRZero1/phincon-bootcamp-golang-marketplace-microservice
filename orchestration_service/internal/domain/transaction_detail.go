package domain

import "time"

type TransactionDetail struct {
	ID            int         `json:"id" binding:"required"`
	TransactionID string      `json:"transaction_id" binding:"required"`
	OrderType     string      `json:"order_type" binding:"required"`
	UserID        int         `json:"user_id" binding:"required"`
	Topic         string      `json:"topic" binding:"required"`
	Action        string      `json:"step" binding:"required"`
	Service       string      `json:"service" binding:"required"`
	Status        string      `json:"status"`
	StatusCode    int         `json:"status_code,omitempty"`
	StatusDesc    string      `json:"status_desc,omitempty"`
	Message       string      `json:"message,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

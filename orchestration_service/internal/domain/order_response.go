package domain

type OrderResponse struct {
	OrderType         string      `json:"order_type"`
	OrderService      string      `json:"order_service"`
	TransactionID     string      `json:"transaction_id"`
	UserID            int         `json:"user_id"`
	Action            string      `json:"action"`
	ResponseCode      int         `json:"response_code"`
	ResponseStatus    string      `json:"response_status"`
	ResponseMessage   string      `json:"response_message"`
	Payload           interface{} `json:"payload"`
	ResponseCreatedAt string      `json:"response_created_at"`
}
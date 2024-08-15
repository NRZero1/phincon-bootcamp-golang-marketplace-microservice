package domain

type IncomingMessage struct {
	OrderType     string `json:"orderType"`
	OrderService  string `json:"orderService,omitempty"` // Omitempty to allow for missing field
	TransactionId string `json:"transactionId"`
	UserId        string `json:"userId"`
	PackageId     string `json:"packageId"`
	RespCode      int    `json:"respCode,omitempty"`    // Omitempty to allow for missing field
	RespStatus    string `json:"respStatus,omitempty"`  // Omitempty to allow for missing field
	RespMessage   string `json:"respMessage,omitempty"` // Omitempty to allow for missing field
}

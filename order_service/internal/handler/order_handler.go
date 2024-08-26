package handler

import "github.com/gin-gonic/gin"

type OrderHandlerInterface interface {
	OrderSave
	// OrderFindTransactionByTransactionID
	// OrderFindTransactionDetailByTransactionID
}

type OrderSave interface {
	Order(context *gin.Context)
}

type OrderFindTransactionByTransactionID interface {
	FindByTransactionID(context *gin.Context)
}

type OrderFindTransactionDetailByTransactionID interface {
	FindTransactionDetailByTransactinID(context *gin.Context)
}

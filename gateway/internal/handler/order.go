package handler

import "github.com/gin-gonic/gin"

type OrderHandlerInterface interface {
}

type OrderSaveTransaction interface {
	SaveTransaction(context *gin.Context)
}

type OrderSaveTransactionDetail interface {
	SaveTransactionDetail(context *gin.Context)
}

type OrderFindByTransactionID interface {
	FindByTransactionID(context *gin.Context)
}


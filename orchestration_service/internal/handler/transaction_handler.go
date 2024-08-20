package handler

import "github.com/gin-gonic/gin"

type TransactionHandlerInterface interface {
	TransactionDetailFindByIDStatusFailed
	TransactionDetailRetry
}

type TransactionDetailFindByIDStatusFailed interface {
	FindTransactionDetailByIDStatusFailed(context *gin.Context)
}

type TransactionDetailRetry interface {
	TransactionDetailRetry(context *gin.Context)
}

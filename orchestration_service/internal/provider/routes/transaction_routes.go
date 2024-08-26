package routes

import (
	"orchestration_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func TransactionRoutes(routerGroup *gin.RouterGroup, transactionHandler handler.TransactionHandlerInterface) {
	routerGroup.GET("/:transaction_id", transactionHandler.FindTransactionDetailByIDStatusFailed)
	routerGroup.PUT("/:id/retry", transactionHandler.TransactionDetailRetry)
}

package routes

import (
	"balance_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func BalanceRoutes(routerGroup *gin.RouterGroup, balanceHandler handler.BalanceHandlerInterface) {
	routerGroup.PATCH("/:id/deduct", balanceHandler.Deduct)
	routerGroup.PATCH("/:id/add", balanceHandler.AddBalance)
	routerGroup.GET("/:id", balanceHandler.FindByID)
	routerGroup.GET("/", balanceHandler.GetAll)
}

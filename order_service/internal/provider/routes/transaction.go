package routes

import (
	"order_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func OrderRoutes(routerGroup *gin.RouterGroup, orderHandler handler.OrderHandlerInterface) {
	// routerGroup.GET("/", orderHandler.GetAll)
	// routerGroup.GET("/:id", orderHandler.FindById)
	routerGroup.POST("/", orderHandler.Order)
}

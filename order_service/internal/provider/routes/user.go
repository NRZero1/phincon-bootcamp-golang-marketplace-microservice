package routes

import (
	"order_service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(routerGroup *gin.RouterGroup, host string) {
	routerGroup.GET("/", middleware.PseudoHandler(host))
	routerGroup.GET("/:id", middleware.PseudoHandler(host))
	routerGroup.POST("/", middleware.PseudoHandler(host))
}

package routes

import (
	"donation_provider_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func DonationRoutes(routerGroup *gin.RouterGroup, balanceHandler handler.DonationProviderHandlerInterface) {
	routerGroup.GET("/:id", balanceHandler.FindByID)
	routerGroup.GET("/", balanceHandler.GetAll)
}

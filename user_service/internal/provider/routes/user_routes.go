package routes

import (
	"user_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func UserRoutes(routerGroup *gin.RouterGroup, userHandler handler.UserHandlerInterface) {
	routerGroup.GET("/", userHandler.GetAllOrFindByName)
	routerGroup.GET("/:id", userHandler.FindById)
	routerGroup.POST("/", userHandler.Save)
	routerGroup.POST("/login", userHandler.Login)
}

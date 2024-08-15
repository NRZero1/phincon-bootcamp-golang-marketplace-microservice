package routes

import (
	"channel_service/internal/handler"

	"github.com/gin-gonic/gin"
)

func ChannelRoutes(routerGroup *gin.RouterGroup, channelHandler handler.ChannelHandlerInterface) {
	routerGroup.GET("/", channelHandler.GetAllOrFindByName)
	routerGroup.GET("/:id", channelHandler.FindById)
}

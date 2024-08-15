package handler

import (
	"channel_service/internal/handler"
	handlerImpl "channel_service/internal/handler/impl"
	providerUseCase "channel_service/internal/provider/usecase"
)

var (
	ChannelHandler handler.ChannelHandlerInterface
)

func InitHandler() {
	ChannelHandler = handlerImpl.NewChannelHandler(providerUseCase.ChannelUseCase)
}

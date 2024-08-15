package handler

import (
	"gateway/internal/handler"
	handlerImpl "gateway/internal/handler/impl"
	providerUseCase "gateway/internal/provider/usecase"
)

var (
	OrderHandler handler.OrderHandlerInterface
)

func InitHandler() {
	OrderHandler = handlerImpl.NewOrderHandler(providerUseCase.OrderUseCase)
}
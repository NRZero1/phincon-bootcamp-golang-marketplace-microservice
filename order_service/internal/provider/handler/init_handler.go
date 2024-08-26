package handler

import (
	"order_service/internal/handler"
	handlerImpl "order_service/internal/handler/impl"
	providerUseCase "order_service/internal/provider/usecase"
)

var (
	OrderHandler handler.OrderHandlerInterface
)

func InitHandler() {
	OrderHandler = handlerImpl.NewOrderHandler(providerUseCase.OrderUseCase)
}

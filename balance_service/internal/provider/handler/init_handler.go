package handler

import (
	"balance_service/internal/handler"
	handlerImpl "balance_service/internal/handler/impl"
	providerUseCase "balance_service/internal/provider/usecase"
)

var (
	BalanceHandler handler.BalanceHandlerInterface
)

func InitHandler() {
	BalanceHandler = handlerImpl.NewBalanceHandler(providerUseCase.BalanceUseCase)
}

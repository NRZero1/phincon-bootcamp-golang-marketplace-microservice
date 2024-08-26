package handler

import (
	"orchestration_service/internal/handler"
	handlerImpl "orchestration_service/internal/handler/impl"
	providerUseCase "orchestration_service/internal/provider/usecase"
)

var (
	TransactionHandler handler.TransactionHandlerInterface
)

func InitHandler() {
	TransactionHandler = handlerImpl.NewTransactionHandler(providerUseCase.TransactionUseCase)
}

package handler

import (
	"user_service/internal/handler"
	handlerImpl "user_service/internal/handler/impl"
	providerUseCase "user_service/internal/provider/usecase"
)

var (
	UserHandler handler.UserHandlerInterface
)

func InitHandler() {
	UserHandler = handlerImpl.NewUserHandler(providerUseCase.UserUseCase)
}

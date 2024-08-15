package usecase

import (
	providerRepo "user_service/internal/provider/repository"
	"user_service/internal/usecase"
	useCaseImpl "user_service/internal/usecase/impl"
)

var (
	UserUseCase usecase.UserUseCaseInterface
)

func InitUseCase() {
	UserUseCase = useCaseImpl.NewUserUseCase(providerRepo.UserRepository)
}

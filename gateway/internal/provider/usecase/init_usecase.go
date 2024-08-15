package usecase

import (
	providerRepo "gateway/internal/provider/repository"
	"gateway/internal/usecase"
	useCaseImpl "gateway/internal/usecase/impl"
)

var (
	OrderUseCase usecase.OrderUseCaseInterface
)

func InitUseCase() {
	OrderUseCase = useCaseImpl.NewOrderUseCase(providerRepo.OrderRepository)
}
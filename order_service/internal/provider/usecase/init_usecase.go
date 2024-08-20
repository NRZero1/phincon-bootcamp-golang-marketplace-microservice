package usecase

import (
	providerRepo "order_service/internal/provider/repository"
	"order_service/internal/usecase"
	useCaseImpl "order_service/internal/usecase/impl"
)

var (
	OrderUseCase usecase.OrderUseCaseInterface
)

func InitUseCase() {
	OrderUseCase = useCaseImpl.NewOrderUseCase(providerRepo.OrderRepository)
}

package usecase

import (
	providerRepo "balance_service/internal/provider/repository"
	"balance_service/internal/usecase"
	useCaseImpl "balance_service/internal/usecase/impl"
)

var (
	BalanceUseCase usecase.BalanceUseCaseInterface
)

func InitUseCase() {
	BalanceUseCase = useCaseImpl.NewBalanceUseCase(providerRepo.BalanceRepository)
}

package usecase

import (
	providerRepo "orchestration_service/internal/provider/repository"
	"orchestration_service/internal/usecase"
	useCaseImpl "orchestration_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
	ConfigUseCase usecase.ConfigUseCaseInterface
	TransactionUseCase usecase.TransactionUseCaseInterface
)

func InitUseCase() {
	ConfigUseCase = useCaseImpl.NewConfigUseCase(providerRepo.ConfigRepository)
	TransactionUseCase = useCaseImpl.NewTransactionUseCase(providerRepo.TransactionRepository)
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(ConfigUseCase, TransactionUseCase)
}

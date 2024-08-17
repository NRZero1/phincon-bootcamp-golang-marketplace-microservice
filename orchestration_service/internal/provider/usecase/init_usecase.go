package usecase

import (
	providerRepo "orchestration_service/internal/provider/repository"
	"orchestration_service/internal/usecase"
	useCaseImpl "orchestration_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
)

func InitUseCase() {
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(providerRepo.ConfigRepository)
}

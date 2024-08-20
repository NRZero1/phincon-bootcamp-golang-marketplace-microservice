package usecase

import (
	"balance_kafka_service/internal/usecase"
	useCaseImpl "balance_kafka_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
	BalanceRequestUseCase usecase.BalanceRequestUseCaseInterface
)

func InitUseCase() {
	BalanceRequestUseCase = useCaseImpl.NewBalanceRequestUseCase()
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(BalanceRequestUseCase)
}

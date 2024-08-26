package usecase

import (
	"user_kafka_service/internal/usecase"
	useCaseImpl "user_kafka_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
	UserRequestUseCase usecase.UserRequestUseCaseInterface
)

func InitUseCase() {
	UserRequestUseCase = useCaseImpl.NewUserRequestUseCase()
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(UserRequestUseCase)
}

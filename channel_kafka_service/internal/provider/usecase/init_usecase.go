package usecase

import (
	"channel_kafka_service/internal/usecase"
	useCaseImpl "channel_kafka_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
	ChannelRequestUseCase usecase.ChannelRequestUseCaseInterface
)

func InitUseCase() {
	ChannelRequestUseCase = useCaseImpl.NewChannelRequestUseCase()
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(ChannelRequestUseCase)
}

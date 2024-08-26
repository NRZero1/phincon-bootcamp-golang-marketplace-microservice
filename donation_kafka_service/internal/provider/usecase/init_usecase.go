package usecase

import (
	"donation_kafka_service/internal/usecase"
	useCaseImpl "donation_kafka_service/internal/usecase/impl"
)

var (
	ConsumerUseCase usecase.ConsumerUseCaseInterface
	DonationRequestUseCase usecase.DonationRequestUseCaseInterface
)

func InitUseCase() {
	DonationRequestUseCase = useCaseImpl.NewDonationRequestUseCase()
	ConsumerUseCase = useCaseImpl.NewConsumerUseCase(DonationRequestUseCase)
}

package usecase

import (
	providerRepo "channel_service/internal/provider/repository"
	"channel_service/internal/usecase"
	useCaseImpl "channel_service/internal/usecase/impl"
)

var (
	ChannelUseCase usecase.ChannelUseCaseInterface
)

func InitUseCase() {
	ChannelUseCase = useCaseImpl.NewChannelUseCase(providerRepo.ChannelRepository)
}

package impl

import (
	"channel_service/internal/domain"
	"channel_service/internal/repository"
	"channel_service/internal/usecase"

	"github.com/rs/zerolog/log"
)

type ChannelUseCase struct {
	repo repository.ChannelRepositoryInterface
}

func NewChannelUseCase(repo repository.ChannelRepositoryInterface) usecase.ChannelUseCaseInterface {
	return ChannelUseCase{
		repo: repo,
	}
}

func (uc ChannelUseCase) FindById(id int) (domain.Channel, error) {
	log.Trace().Msg("Entering channel usecase find by id")
	return uc.repo.FindById(id)
}

func (uc ChannelUseCase) GetAll() []domain.Channel {
	log.Trace().Msg("Entering channel usecase get all")
	return uc.repo.GetAll()
}

func (uc ChannelUseCase) FindByName(name string) (domain.Channel, error) {
	log.Trace().Msg("Entering channel usecase FindByName")
	return uc.repo.FindByName(name)
}

func (uc ChannelUseCase) AddMembership(channelID int, userID int) (domain.Channel, error) {
	log.Trace().Msg("Entering channel Add Membership")
	return uc.repo.AddMembership(channelID, userID)
}

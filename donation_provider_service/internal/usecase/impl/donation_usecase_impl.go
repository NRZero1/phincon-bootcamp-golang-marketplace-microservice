package impl

import (
	"donation_provider_service/internal/domain"
	"donation_provider_service/internal/repository"
	"donation_provider_service/internal/usecase"

	"github.com/rs/zerolog/log"
)

type DonationProviderUseCase struct {
	repo repository.DonationProviderRepositoryInterface
}

func NewDonationProviderUseCase(repo repository.DonationProviderRepositoryInterface) usecase.DonationProviderUseCaseInterface {
	return DonationProviderUseCase{
		repo: repo,
	}
}

func (uc DonationProviderUseCase) FindByID(id int) (domain.Provider, error) {
	log.Trace().Msg("Entering donation_provider usecase find by id")
	return uc.repo.FindByID(id)
}

func (uc DonationProviderUseCase) GetAll() []domain.Provider {
	log.Trace().Msg("Entering donation_provider usecase get all")
	return uc.repo.GetAll()
}

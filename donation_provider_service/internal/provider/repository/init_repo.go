package repository

import (
	"donation_provider_service/internal/repository"
	repoImplement "donation_provider_service/internal/repository/impl"
)

var (
	DonationProviderRepository repository.DonationProviderRepositoryInterface
)

func InitRepository() {
	DonationProviderRepository = repoImplement.NewDonationProviderRepository()
}

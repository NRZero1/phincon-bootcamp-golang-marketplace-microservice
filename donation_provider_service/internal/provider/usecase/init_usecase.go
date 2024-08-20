package usecase

import (
	providerRepo "donation_provider_service/internal/provider/repository"
	"donation_provider_service/internal/usecase"
	useCaseImpl "donation_provider_service/internal/usecase/impl"
)

var (
	DonationProviderUseCase usecase.DonationProviderUseCaseInterface
)

func InitUseCase() {
	DonationProviderUseCase = useCaseImpl.NewDonationProviderUseCase(providerRepo.DonationProviderRepository)
}

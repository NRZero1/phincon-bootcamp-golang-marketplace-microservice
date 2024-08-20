package usecase

import "donation_provider_service/internal/domain"

type DonationProviderUseCaseInterface interface {
	DonationProviderFindByID
	DonationProviderGetAll
}

type DonationProviderFindByID interface {
	FindByID(id int) (domain.Provider, error)
}

type DonationProviderGetAll interface {
	GetAll() ([]domain.Provider)
}

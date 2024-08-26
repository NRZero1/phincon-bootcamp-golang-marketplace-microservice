package repository

import (
	"donation_provider_service/internal/domain"
)

type DonationProviderRepositoryInterface interface {
	DonationProviderFindByID
	DonationProviderGetAll
}

type DonationProviderFindByID interface {
	FindByID(id int) (domain.Provider, error)
}

type DonationProviderGetAll interface {
	GetAll() ([]domain.Provider)
}

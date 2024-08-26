package usecase

type DonationRequestUseCaseInterface interface {
	ProviderGetByID
}

type ProviderGetByID interface {
	GetProviderByID(id int) (bool, int, error)
}

package handler

import (
	"donation_provider_service/internal/handler"
	handlerImpl "donation_provider_service/internal/handler/impl"
	providerUseCase "donation_provider_service/internal/provider/usecase"
)

var (
	DonationProviderHandler handler.DonationProviderHandlerInterface
)

func InitHandler() {
	DonationProviderHandler = handlerImpl.NewBalanceHandler(providerUseCase.DonationProviderUseCase)
}

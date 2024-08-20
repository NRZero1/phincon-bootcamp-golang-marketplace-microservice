package handler

import "github.com/gin-gonic/gin"

type DonationProviderHandlerInterface interface {
	DonationProviderFindByID
	DonationProviderGetAll
}

type DonationProviderFindByID interface {
	FindByID(context *gin.Context)
}

type DonationProviderGetAll interface {
	GetAll(context *gin.Context)
}

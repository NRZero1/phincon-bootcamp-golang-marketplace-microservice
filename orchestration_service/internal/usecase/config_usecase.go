package usecase

import (
	"orchestration_service/internal/domain"
)

type ConfigUseCaseInterface interface {
	ConfigGetConfigByOrderType
}

type ConfigGetConfigByOrderType interface {
	GetConfigByOrderType(orderType string, serviceSource string, statusCategory string) (domain.ConfigResponse, error)
}

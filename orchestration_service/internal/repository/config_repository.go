package repository

import (
	"context"
	"orchestration_service/internal/domain"
)

type ConfigRepositoryInterface interface {
	ConfigGetConfigByOrderType
}

type ConfigGetConfigByOrderType interface {
	GetConfigByOrderType(context context.Context, orderType string, serviceSource string, statusCategory string) (domain.ConfigResponse, error)
}

package impl

import (
	"context"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/usecase"
	"time"

	"github.com/rs/zerolog/log"
)

type ConfigUseCase struct {
	repo repository.ConfigRepositoryInterface
}

func NewConfigUseCase(repo repository.ConfigRepositoryInterface) usecase.ConfigUseCaseInterface {
	return ConfigUseCase {
		repo: repo,
	}
}

func (uc ConfigUseCase) GetConfigByOrderType(orderType string, serviceSource string, statusCategory string) (domain.ConfigResponse, error) {
	log.Trace().Msg("Inside Get Config usecase ")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.GetConfigByOrderType(ctx, orderType, serviceSource, statusCategory)
}

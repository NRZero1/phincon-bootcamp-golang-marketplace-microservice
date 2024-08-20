package impl

import (
	"balance_service/internal/domain"
	"balance_service/internal/repository"
	"balance_service/internal/usecase"

	"github.com/rs/zerolog/log"
)

type BalanceUseCase struct {
	repo repository.BalanceRepositoryInterface
}

func NewBalanceUseCase(repo repository.BalanceRepositoryInterface) usecase.BalanceUseCaseInterface {
	return BalanceUseCase{
		repo: repo,
	}
}

func (uc BalanceUseCase) FindByID(id int) (domain.Balance, error) {
	log.Trace().Msg("Entering balance usecase find by id")
	return uc.repo.FindByID(id)
}

func (uc BalanceUseCase) GetAll() []domain.Balance {
	log.Trace().Msg("Entering balance usecase get all")
	return uc.repo.GetAll()
}

func (uc BalanceUseCase) Deduct(userID int, amount float64) (error) {
	log.Trace().Msg("Entering balance deduct use case")
	return uc.repo.Deduct(userID, amount)
}

func (uc BalanceUseCase) AddBalance(userID int, amount float64) (error) {
	log.Trace().Msg("Entering add balance use case")
	return uc.repo.AddBalance(userID, amount)
}

package repository

import (
	"balance_service/internal/repository"
	repoImplement "balance_service/internal/repository/impl"
)

var (
	BalanceRepository repository.BalanceRepositoryInterface
)

func InitRepository() {
	BalanceRepository = repoImplement.NewBalanceRepository()
}

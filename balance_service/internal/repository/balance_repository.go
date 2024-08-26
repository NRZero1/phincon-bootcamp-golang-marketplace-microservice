package repository

import (
	"balance_service/internal/domain"
)

type BalanceRepositoryInterface interface {
	BalanceFindByID
	BalanceGetAll
	BalanceDeduct
	BalanceAddBalance
}

type BalanceFindByID interface {
	FindByID(id int) (domain.Balance, error)
}

type BalanceDeduct interface {
	Deduct(userID int, amount float64) (error)
}

type BalanceGetAll interface {
	GetAll() ([]domain.Balance)
}

type BalanceAddBalance interface {
	AddBalance(userID int, amount float64) (error)
}

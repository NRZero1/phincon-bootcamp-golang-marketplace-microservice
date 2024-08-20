package usecase

import "balance_service/internal/domain"

type BalanceUseCaseInterface interface {
	BalanceFindByID
	BalanceGetAll
	BalanceDeduct
	BalanceAddBalance
}

type BalanceFindByID interface {
	FindByID(id int) (domain.Balance, error)
}

type BalanceGetAll interface {
	GetAll() ([]domain.Balance)
}

type BalanceDeduct interface {
	Deduct(id int, amount float64) (error)
}

type BalanceAddBalance interface {
	AddBalance(id int, amount float64) (error)
}

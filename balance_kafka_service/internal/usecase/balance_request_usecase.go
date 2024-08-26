package usecase

type BalanceRequestUseCaseInterface interface {
	BalanceGetUserByID
}

type BalanceGetUserByID interface {
	GetBalanceByID(id int) (bool, int, error)
}

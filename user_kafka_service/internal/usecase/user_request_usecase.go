package usecase

type UserRequestUseCaseInterface interface {
	UserGetUserByID
}

type UserGetUserByID interface {
	GetUserByID(id int) (bool, int, error)
}

package usecase

import (
	"user_service/internal/domain/dto/request"
	"user_service/internal/domain/dto/response"
)

type UserUseCaseInterface interface {
	UserSave
	UserFindById
	UserGetAll
	UserBalanceReduce
	UserFindByUsernameLogin
	UserFindByUsername
	UserSetPackage
}

type UserSave interface {
	Save(user request.Register) (response.UserResponse, error)
}

type UserFindById interface {
	FindById(id int) (response.UserResponse, error)
}

type UserGetAll interface {
	GetAll() ([]response.UserResponse)
}

type UserBalanceReduce interface {
	ReduceBalance(id int, amount float64) (response.UserResponse, error)
}

type UserFindByUsernameLogin interface {
	FindByUsernameLogin(username string) (response.LoginResponse, error)
}

type UserFindByUsername interface {
	FindByUsername(username string) (response.UserResponse, error)
}

type UserSetPackage interface {
	SetPackage(userID int, packageID int) (response.UserResponse, error)
}

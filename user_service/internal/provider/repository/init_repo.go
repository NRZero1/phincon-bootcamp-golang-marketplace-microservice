package repository

import (
	"user_service/internal/repository"
	repoImplement "user_service/internal/repository/impl"
)

var (
	UserRepository repository.UserRepositoryInterface
)

func InitRepository() {
	UserRepository = repoImplement.NewUserRepository()
}

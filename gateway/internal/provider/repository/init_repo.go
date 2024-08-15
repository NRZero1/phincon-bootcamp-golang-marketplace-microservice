package repository

import (
	"database/sql"
	"gateway/internal/repository"
	repoImplement "gateway/internal/repository/impl"
)

var (
	OrderRepository repository.OrderRepositoryInterface
)

func InitRepository(database *sql.DB) {
	OrderRepository = repoImplement.NewOrderRepository(database)
}
package repository

import (
	"database/sql"
	"order_service/internal/repository"
	repoImplement "order_service/internal/repository/impl"
)

var (
	OrderRepository repository.OrderRepositoryInterface
)

func InitRepository(database *sql.DB) {
	OrderRepository = repoImplement.NewOrderRepository(database)
}

package repository

import (
	"database/sql"
	"orchestration_service/internal/repository"
	repoImplement "orchestration_service/internal/repository/impl"
)

var (
	ConfigRepository repository.ConfigRepositoryInterface
)

func InitRepository(database *sql.DB) {
	ConfigRepository = repoImplement.NewConfigRepository(database)
}

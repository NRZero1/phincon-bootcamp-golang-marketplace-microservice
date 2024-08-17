package impl

import (
	"context"
	"database/sql"
	"fmt"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/utils"

	"github.com/rs/zerolog/log"
)

type ConfigRepository struct {
	db *sql.DB
}

func NewConfigRepository(database *sql.DB) repository.ConfigRepositoryInterface {
	return &ConfigRepository {
		db: database,
	}
}

func (repo ConfigRepository) GetConfigByOrderType(ctx context.Context, orderType string, serviceSource string, statusCategory string) (domain.ConfigResponse, error) {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return domain.ConfigResponse{}, utils.ErrRepoCreateTrx
	}

	query := `
		SELECT
			order_type,
			service_source,
			service_dest,
			action,
			status_category
		FROM
			transaction_config
		WHERE
			order_type=$1
		AND
			service_source=$2
		AND
			status_category=$3
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return domain.ConfigResponse{}, utils.ErrPreparedStmt
	}

	defer stmt.Close()

	var configResponse domain.ConfigResponse

	errScan := stmt.QueryRowContext(ctx, orderType).Scan(&configResponse.OrderType, &configResponse.ServiceSource, &configResponse.ServiceDest, &configResponse.Action, &configResponse.StatusCategory)

	if errScan != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to scan config steps with error message: %s", errScan.Error()))
		if errScan == sql.ErrNoRows {
			log.Error().Msg(fmt.Sprintf("Config with order type %s not found", orderType))
			return domain.ConfigResponse{}, utils.ErrNoSqlRows
		}
		return domain.ConfigResponse{}, utils.ErrErrScan
	}

	if err = trx.Commit(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to commit trx with message: %s", err.Error()))
		return domain.ConfigResponse{}, err
	}

	log.Debug().Msg(fmt.Sprintf("Config fetched with value %+v", configResponse))
	return configResponse, nil
}

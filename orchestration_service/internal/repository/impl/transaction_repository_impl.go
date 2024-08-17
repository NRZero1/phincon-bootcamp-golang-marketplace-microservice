package impl

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"orchestration_service/internal/domain"
// 	"orchestration_service/internal/repository"
// 	"orchestration_service/internal/utils"

// 	"github.com/rs/zerolog/log"
// )

// type TranscationRepository struct {
// 	db *sql.DB
// }

// func NewTranscationRepository(database *sql.DB) repository.TransactionRepositoryInterface {
// 	return &TranscationRepository{
// 		db: database,
// 	}
// }

// func (repo TranscationRepository) TransactionUpdate(ctx context.Context, transaction domain.Transaction) (error) {
// 	trx, err := repo.db.BeginTx(ctx, nil)

// 	if err != nil {
// 		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
// 		return utils.ErrRepoCreateTrx
// 	}

// 	query := `
// 		UPDATE
// 			transaction
// 		SET
// 			status=$1
// 			update_at=$2
// 		WHERE
// 			transaction_id=$3
// 	`

// 	stmt, err := trx.PrepareContext(ctx, query)

// 	if err != nil {
// 		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
// 		return utils.ErrPreparedStmt
// 	}

// 	defer stmt.Close()

// 	_, errExec := stmt.ExecContext(ctx, , transaction.OrderType, transaction.UserID, "CREATED")

// 	if errExec != nil {
// 		trx.Rollback()
// 		log.Error().Msg(fmt.Sprintf("Error when trying to exec statement in repo with message: %s", errExec.Error()))
// 		return utils.ErrDbExec
// 	}

// 	if err = trx.Commit(); err != nil {
// 		log.Error().Msg(fmt.Sprintf("Error when trying to commit trx in repo with message: %s", err.Error()))
// 		return utils.ErrTrxCommit
// 	}

// 	return nil
// }

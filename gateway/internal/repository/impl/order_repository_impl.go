package impl

import (
	"context"
	"database/sql"
	"fmt"
	"gateway/internal/domain"
	"gateway/internal/repository"
	"gateway/internal/utils"

	"github.com/rs/zerolog/log"
)

type OrderRepository struct {
	db *sql.DB
}

func NewOrderRepository(database *sql.DB) repository.OrderRepositoryInterface {
	return &OrderRepository {
		db: database,
	}
}

func (repo *OrderRepository) SaveTransaction(ctx context.Context, transaction domain.Transaction) error {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	query := "INSERT INTO transaction (transaction_id, order_type, user_id, status) VALUES ($1, $2, $3, $4)"

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	_, errExec := stmt.ExecContext(ctx, transaction.TransactionID, transaction.OrderType, transaction.UserID, "CREATED")

	if errExec != nil {
		trx.Rollback()
		log.Error().Msg(fmt.Sprintf("Error when trying to exec statement in repo with message: %s", errExec.Error()))
		return utils.ErrDbExec
	}

	if err = trx.Commit(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to commit trx in repo with message: %s", err.Error()))
		return utils.ErrTrxCommit
	}

	return nil
}

func (repo *OrderRepository) GetLastTransactionID(ctx context.Context) (string, error) {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return "", utils.ErrRepoCreateTrx
	}

	query := `
		WITH latest_transaction AS (
			SELECT transaction_id
			FROM transaction
			WHERE transaction_id LIKE CONCAT(to_char(current_date, 'YYMMDD'), '%')
			ORDER BY transaction_id DESC
			LIMIT 1
		)
		SELECT transaction_id
		FROM latest_transaction;
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return "", utils.ErrPreparedStmt
	}

	defer stmt.Close()

	var transactionID string

	errScan := stmt.QueryRowContext(ctx).Scan(&transactionID)

	if errScan != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to scan query result in repo with message: %s", errScan.Error()))
		if errScan == sql.ErrNoRows {
			return "", utils.ErrNoSqlRows
		}
		return "", utils.ErrErrScan
	}

	return transactionID, nil
}

func (repo *OrderRepository) SaveTransactionDetails(ctx context.Context, transactionDetail domain.TransactionDetail) error {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	query := `
	INSERT INTO
		transaction_detail
		(
			transaction_id,
			order_type,
			user_id,
			topic,
			action,
			service,
			status,
			status_code,
			status_desc,
			message,
			payload
		)
		VALUES
		(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11
		)
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	_, errExec := stmt.ExecContext(
		ctx,
		transactionDetail.TransactionID,
		transactionDetail.OrderType,
		transactionDetail.UserID,
		transactionDetail.Topic,
		transactionDetail.Action,
		transactionDetail.Service,
		transactionDetail.Status,
		transactionDetail.StatusCode,
		transactionDetail.StatusDesc,
		transactionDetail.Message,
		transactionDetail.Payload,
	)

	if errExec != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to exec statement in repo with message: %s", errExec.Error()))
		trx.Rollback()
		return utils.ErrDbExec
	}

	if err = trx.Commit(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to commit trx in repo with message: %s", err.Error()))
		return utils.ErrTrxCommit
	}

	return nil
}

func (repo *OrderRepository) FindByTransactionID(ctx context.Context, transactionID string) (domain.Transaction, error) {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return domain.Transaction{}, utils.ErrRepoCreateTrx
	}

	query := `
		SELECT
			id,
			transaction_id,
			order_type,
			user_id,
			status,
			created_at,
			updated_at
		FROM
			transaction
		WHERE
			transaction_id=$1
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return domain.Transaction{}, utils.ErrPreparedStmt
	}

	defer stmt.Close()

	var transaction domain.Transaction
	errScan := stmt.QueryRowContext(ctx, transactionID).
		Scan(
			&transaction.ID,
			&transaction.TransactionID,
			&transaction.OrderType,
			&transaction.UserID,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

	if errScan != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to scan query result in repo with message: %s", errScan.Error()))
		if errScan == sql.ErrNoRows {
			return domain.Transaction{}, utils.ErrNoSqlRows
		}
		return domain.Transaction{}, utils.ErrErrScan
	}

	return transaction, nil
}

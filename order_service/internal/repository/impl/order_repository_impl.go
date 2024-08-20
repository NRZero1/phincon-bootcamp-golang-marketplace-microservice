package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"order_service/internal/domain"
	"order_service/internal/repository"
	"order_service/internal/utils"

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
	log.Trace().Msg("Inside order repo to save transaction")
	log.Trace().Msg("Begin trx")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Setting insert query")
	query := "INSERT INTO transaction (transaction_id, order_type, user_id, status) VALUES ($1, $2, $3, $4)"

	log.Trace().Msg("Trying to create prepared statement")
	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	log.Trace().Msg("Trying to exec query")
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

	log.Trace().Msg("Order repo save transaction completed")
	return nil
}

func (repo *OrderRepository) GetLastTransactionID(ctx context.Context) (string, error) {
	log.Trace().Msg("Begin trx")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return "", utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Setting query")
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

	log.Trace().Msg("Preparing prepared statement")
	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return "", utils.ErrPreparedStmt
	}

	defer stmt.Close()

	var transactionID string

	log.Trace().Msg("Trying to scan result from query row context")
	errScan := stmt.QueryRowContext(ctx).Scan(&transactionID)

	if errScan != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to scan query result in repo with message: %s", errScan.Error()))
		if errScan == sql.ErrNoRows {
			return "", nil
		}
		return "", utils.ErrErrScan
	}

	return transactionID, nil
}

func (repo *OrderRepository) SaveTransactionDetails(ctx context.Context, transactionDetail domain.TransactionDetail) error {
	log.Trace().Msg("Inside order repo save transaction detail")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Setting query")
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

	log.Trace().Msg("Trying to create prepared statement")
	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	payloadJSON, err := json.Marshal(transactionDetail.Payload)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to marshal payload in repo with message: %s", err.Error()))
		trx.Rollback()
		return utils.ErrMarshalPayload
	}

	log.Trace().Msg("Trying to exec query")
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
		string(payloadJSON),
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

	log.Trace().Msg("Saving transaction detail completed")
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

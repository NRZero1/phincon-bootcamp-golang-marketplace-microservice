package impl

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/utils"
	"time"

	"github.com/rs/zerolog/log"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(database *sql.DB) repository.TransactionRepositoryInterface {
	return &TransactionRepository{
		db: database,
	}
}

func (repo TransactionRepository) TransactionUpdate(ctx context.Context, status string, transactionID string) (error) {
	log.Trace().Msg("Inside transaction repo transaction update")
	log.Trace().Msg("Begin trx")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Set up query")
	query := `
		UPDATE
			transaction
		SET
			status=$1,
			update_at=$2
		WHERE
			transaction_id=$3
	`

	log.Trace().Msg("Trying to create prepared statement")
	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	log.Debug().Str("status: ", status).Msg("Debug")
	_, errExec := stmt.ExecContext(ctx, status, time.Now(), transactionID)

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

func (repo TransactionRepository) TransactionDetailInput(ctx context.Context, transactionMessage domain.TransactionMessage, payload interface{}) error {
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	query := `
		INSERT INTO transaction_detail
		(transaction_id, order_type, user_id, topic, action, service, status, status_code, status_desc, message, payload)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	_, errExec := stmt.ExecContext(
		ctx,
		transactionMessage.TransactionID,
		transactionMessage.OrderType,
		transactionMessage.UserID,
		transactionMessage.Topic,
		transactionMessage.Action,
		transactionMessage.Service,
		transactionMessage.Status,
		transactionMessage.StatusCode,
		transactionMessage.StatusDesc,
		transactionMessage.Message,
		payload,
	)

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

func (repo TransactionRepository) FindTransactionDetailByIDStatusFailed(ctx context.Context, transactionID string) (domain.TransactionDetail, error) {
	log.Trace().Msg("Inside repo FindTransactionDetailByIDStatusFailed")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Set query")
	query := `
		SELECT
			id,
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
			payload,
			created_at
		FROM
			transaction_detail
		WHERE
			transaction_id=$1
		AND
			status='FAILED'
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrPreparedStmt
	}

	defer stmt.Close()

	var transactionDetail domain.TransactionDetail
	errScan := stmt.QueryRowContext(ctx, transactionID).Scan(
		&transactionDetail.ID,
		&transactionDetail.TransactionID,
		&transactionDetail.OrderType,
		&transactionDetail.UserID,
		&transactionDetail.Topic,
		&transactionDetail.Action,
		&transactionDetail.Service,
		&transactionDetail.Status,
		&transactionDetail.StatusCode,
		&transactionDetail.StatusDesc,
		&transactionDetail.Message,
		&transactionDetail.Payload,
		&transactionDetail.CreatedAt,
	)

	if errScan != nil {
		if errScan == sql.ErrNoRows {
			log.Error().Msg(fmt.Sprintf("Transaction Detail with TransactionID %s not found", transactionID))
			return domain.TransactionDetail{}, utils.ErrNoSqlRows
		}
		return domain.TransactionDetail{}, utils.ErrErrScan
	}

	var decodedPayload []byte
	switch p := transactionDetail.Payload.(type) {
	case string:
		decodedPayload, err = base64.StdEncoding.DecodeString(p)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when decoding base64 payload: %s", err.Error()))
			return domain.TransactionDetail{}, utils.ErrInvalidPayloadFormat
		}
		err = json.Unmarshal(decodedPayload, &transactionDetail.Payload)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when unmarshaling payload: %s", err.Error()))
			return domain.TransactionDetail{}, utils.ErrInvalidPayloadFormat
		}
	case []byte:
		decodedPayload = p
		err = json.Unmarshal(decodedPayload, &transactionDetail.Payload)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when unmarshaling payload: %s", err.Error()))
			return domain.TransactionDetail{}, utils.ErrInvalidPayloadFormat
		}
	default:
		log.Error().Msg("Unexpected type for Payload")
		return domain.TransactionDetail{}, utils.ErrInvalidPayloadFormat
	}

	if err = trx.Commit(); err != nil {
		return domain.TransactionDetail{}, utils.ErrTrxCommit
	}

	log.Info().Msg("Find transaction detail by ID completed")
	return transactionDetail, nil
}

func (repo TransactionRepository) TransactionDetailRetry(ctx context.Context, id int, transactionDetail domain.TransactionDetail) (domain.TransactionDetail, error) {
	log.Trace().Msg("Inside transaction repo transaction update")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Set up query")
	query := `
		UPDATE
			transaction_detail
		SET
			transaction_id=$1,
			order_type=$2,
			user_id=$3,
			topic=$4,
			action=$5,
			service=$6,
			status=$7,
			status_code=$8,
			status_desc=$9,
			message=$10,
			payload=$11,
			created_at=$12
		WHERE
			id=$13
		RETURNING
			id,
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
			payload,
			created_at
	`

	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrPreparedStmt
	}

	defer stmt.Close()

	// Marshal Payload to JSON
	payloadBytes, err := json.Marshal(transactionDetail.Payload)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when marshaling Payload to JSON: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrInvalidPayloadFormat
	}

	// Convert payload to string
	payloadJSON := string(payloadBytes)
	log.Debug().Msgf("Payload before insert/update: %s", payloadJSON)

	var updatedTransactionDetail domain.TransactionDetail

	errScan := stmt.QueryRowContext(
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
		payloadJSON, // Pass as JSON string directly
		transactionDetail.CreatedAt,
		transactionDetail.ID,
	).Scan(
		&updatedTransactionDetail.ID,
		&updatedTransactionDetail.TransactionID,
		&updatedTransactionDetail.OrderType,
		&updatedTransactionDetail.UserID,
		&updatedTransactionDetail.Topic,
		&updatedTransactionDetail.Action,
		&updatedTransactionDetail.Service,
		&updatedTransactionDetail.Status,
		&updatedTransactionDetail.StatusCode,
		&updatedTransactionDetail.StatusDesc,
		&updatedTransactionDetail.Message,
		&updatedTransactionDetail.Payload,
		&updatedTransactionDetail.CreatedAt,
	)

	if errScan != nil {
		trx.Rollback()
		log.Error().Msg(fmt.Sprintf("Error when trying to execute statement in repo with message: %s", errScan.Error()))
		return domain.TransactionDetail{}, utils.ErrErrScan
	}

	if err = trx.Commit(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to commit trx in repo with message: %s", err.Error()))
		return domain.TransactionDetail{}, utils.ErrTrxCommit
	}

	return updatedTransactionDetail, nil
}

func (repo TransactionRepository) UpdateStatus(ctx context.Context, status string, transactionID string) error {
	log.Trace().Msg("AAAAAAAAAAAAAAAAAAAAAAAAAAA")
	log.Trace().Msg("Begin trx")
	trx, err := repo.db.BeginTx(ctx, nil)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create trx in repo with message: %s", err.Error()))
		return utils.ErrRepoCreateTrx
	}

	log.Trace().Msg("Set up query")
	query := `
		UPDATE
			transaction_detail
		SET
			status=$1
		WHERE
			transaction_id=$2
	`

	log.Trace().Msg("Trying to create prepared statement")
	stmt, err := trx.PrepareContext(ctx, query)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create prepared statement in repo with message: %s", err.Error()))
		return utils.ErrPreparedStmt
	}

	defer stmt.Close()

	log.Debug().Msgf("TransactionID: %s", transactionID)
	log.Debug().Msgf("Status: %s", status)
	_, errExec := stmt.ExecContext(ctx, status, transactionID)

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

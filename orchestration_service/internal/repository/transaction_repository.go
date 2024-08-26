package repository

import (
	"context"
	"orchestration_service/internal/domain"
)

type TransactionRepositoryInterface interface {
	TransactionUpdate
	TransactionDetailInput
	TransactionDetailFindByIDStatusFailed
	TransactionDetailRetry
	TransactionDetailUpdateStatus
}

type TransactionUpdate interface {
	TransactionUpdate(context context.Context, status string, transactionID string) (error)
}

type TransactionDetailInput interface {
	TransactionDetailInput(context context.Context, transactionMessage domain.TransactionMessage, payload interface{}) (error)
}

type TransactionDetailFindByIDStatusFailed interface {
	FindTransactionDetailByIDStatusFailed(context context.Context, transactionID string) (domain.TransactionDetail, error)
}

type TransactionDetailRetry interface {
	TransactionDetailRetry(context context.Context, id int, transactionDetail domain.TransactionDetail) (domain.TransactionDetail, error)
}

type TransactionDetailUpdateStatus interface {
	UpdateStatus(context context.Context, status string, transactionID string) error
}

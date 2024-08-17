package repository

import (
	"context"
	"orchestration_service/internal/domain"
)

type TransactionRepositoryInterface interface {
}

type TransactionUpdate interface {
	TransactionUpdate(context context.Context, transaction domain.Transaction) (error)
}

type TransactionDetailUpdate interface {
	TransactionDetailUpdate(context context.Context, transactionDetail domain.TransactionDetail) (error)
}

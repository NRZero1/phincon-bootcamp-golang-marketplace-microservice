package repository

import (
	"context"
	"gateway/internal/domain"
)

type OrderRepositoryInterface interface {
	OrderGetLastTransactionID
	OrderFindByTransactionID
	OrderSaveTransaction
	OrderSaveTransactionDetails
}

type OrderGetLastTransactionID interface {
	GetLastTransactionID(context context.Context) (string, error)
}

type OrderSaveTransaction interface {
	SaveTransaction(ctx context.Context, transaction domain.Transaction) error
}

type OrderSaveTransactionDetails interface {
	SaveTransactionDetails(ctx context.Context, transactionDetail domain.TransactionDetail) error
}

type OrderFindByTransactionID interface {
	FindByTransactionID(ctx context.Context, transactionID string) (domain.Transaction, error)
}
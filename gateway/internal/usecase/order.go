package usecase

import (
	"gateway/internal/domain"
)

type OrderUseCaseInterface interface {
	OrderGetLastTransactionID
	OrderFindByTransactionID
	OrderSaveTransaction
	OrderSaveTransactionDetails
}

type OrderGetLastTransactionID interface {
	GetLastTransactionID() (string, error)
}

type OrderSaveTransaction interface {
	SaveTransaction(orderRequest domain.OrderRequest) error
}

type OrderSaveTransactionDetails interface {
	SaveTransactionDetails(transactionDetail domain.TransactionDetail) error
}

type OrderFindByTransactionID interface {
	FindByTransactionID(transactionID string) (domain.Transaction, error)
}
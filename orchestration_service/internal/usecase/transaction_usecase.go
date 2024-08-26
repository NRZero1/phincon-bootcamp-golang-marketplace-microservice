package usecase

import (
	"orchestration_service/internal/domain"
)

type TransactionUseCaseInterface interface {
	TransactionUpdate
	TransactionDetailInput
	TransactionDetailFindByIDStatusFailed
	TransactionDetailRetry
	TransactionDetailSend
	TransactionDetailUpdateStatus
}

type TransactionUpdate interface {
	TransactionUpdate(status string, transactionID string) error
}

type TransactionDetailInput interface {
	TransactionDetailInput(transactionMessage domain.TransactionMessage, payload interface{}) error
}
type TransactionDetailFindByIDStatusFailed interface {
	FindTransactionDetailByIDStatusFailed(transactionID string) (domain.TransactionDetail, error)
}

type TransactionDetailRetry interface {
	TransactionDetailRetry(id int, transactionDetail domain.TransactionDetail) (domain.TransactionDetail, error)
}

type TransactionDetailSend interface {
	TransactionDetailSend(transactionDetail domain.TransactionDetail)
}

type TransactionDetailUpdateStatus interface {
	UpdateStatus(status string, transactionID string) error
}


package impl

import (
	"context"
	"gateway/internal/domain"
	"gateway/internal/provider/kafka"
	"gateway/internal/repository"
	"gateway/internal/usecase"
	"gateway/internal/utils"
	"time"
)

type OrderUseCase struct {
	repo repository.OrderRepositoryInterface
}

func NewOrderUseCase(repo repository.OrderRepositoryInterface) usecase.OrderUseCaseInterface {
	return OrderUseCase {
		repo: repo,
	}
}

func (uc OrderUseCase) SaveTransaction(orderRequest domain.OrderRequest) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	lastID, err := uc.GetLastTransactionID()

	if err != nil {
		return err
	}
	generatedID, err := utils.GenerateNextTransactionID(lastID)

	if err != nil {
		return err
	}
	
	transaction := domain.Transaction {
		TransactionID: generatedID,
		OrderType: orderRequest.OrderType,
		Status: "CREATED",
		UserID: orderRequest.UserID,
	}

	transactionDetail := domain.TransactionDetail {
		TransactionID: generatedID,
		OrderType: orderRequest.OrderType,
		Topic: "orchestration",
		Step: "start",
		Service: "gateway",
		Status: "ONPROGRESS",
		Message: "Sending message to orchestration service from gateway service",
		Payload: orderRequest.Payload,
	}

	errSave := uc.repo.SaveTransaction(ctx, transaction)

	if errSave != nil {
		return errSave
	}

	errSaveDetail := uc.SaveTransactionDetails(transactionDetail)

	kafka.NewKafkaProducer()
	
	return errSaveDetail
}

func (uc OrderUseCase) GetLastTransactionID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)

	defer cancel()
	lastID, err := uc.repo.GetLastTransactionID(ctx)

	return lastID, err
}

func (uc OrderUseCase) SaveTransactionDetails(transactionDetail domain.TransactionDetail) (error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	err := uc.repo.SaveTransactionDetails(ctx, transactionDetail)

	return err
}

func (uc OrderUseCase) FindByTransactionID(transactionID string) (domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.FindByTransactionID(ctx, transactionID)
}
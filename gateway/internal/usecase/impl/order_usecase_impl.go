package impl

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	"gateway/internal/domain/payload/request"
	"gateway/internal/provider/kafka"
	"gateway/internal/repository"
	"gateway/internal/usecase"
	"gateway/internal/utils"
	"net/http"
	"os"
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

func (uc OrderUseCase) SaveTransaction(orderRequest domain.OrderRequest) (domain.OrderResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	lastID, err := uc.GetLastTransactionID()

	if err != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
	}
	generatedID, err := utils.GenerateNextTransactionID(lastID)

	if err != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
	}

	orderRequest.TransactionID = generatedID

	transaction := domain.Transaction {
		TransactionID: generatedID,
		OrderType: orderRequest.OrderType,
		Status: "CREATED",
		UserID: orderRequest.UserID,
	}

	var transactionDetail domain.TransactionDetail

	if orderRequest.OrderType == "membership" {
		transactionDetail = domain.TransactionDetail {
			TransactionID: generatedID,
			OrderType: orderRequest.OrderType,
			Topic: "orchestration",
			Action: "START",
			Service: "gateway",
			Status: "ONPROGRESS",
			Payload: request.MembershipRequest {
				ChannelID: orderRequest.Payload.(request.MembershipRequest).ChannelID,
			},
			UserID: orderRequest.UserID,
		}
	}

	errSave := uc.repo.SaveTransaction(ctx, transaction)

	if errSave != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, errSave.Error()), errSave
	}

	errSaveDetail := uc.SaveTransactionDetails(transactionDetail)

	if errSaveDetail != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, fmt.Sprintf("Transaction created but couldn't write the details with message: %s", errSaveDetail.Error())), errSaveDetail
	}

	broker := []string{os.Getenv("KAFKA_BROKER")}
	producer, err := kafka.NewKafkaProducer(broker, transactionDetail.Topic, 1, 1)

	if err != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
	}

	var message any

	if transactionDetail.OrderType == "membership" {
		message = domain.Message[domain.TransactionMessage, request.MembershipRequest] {
			Header: domain.TransactionMessage{
				TransactionID: transactionDetail.TransactionID,
				OrderType: transactionDetail.OrderType,
				UserID: transactionDetail.UserID,
				Topic: transactionDetail.Topic,
				Action: transactionDetail.Topic,
				Service: transactionDetail.Service,
				Status: transactionDetail.Status,
			},
			Body: request.MembershipRequest{
				ChannelID: orderRequest.Payload.(request.MembershipRequest).ChannelID,
			},
		}
	}

	errProduceMessage := producer.ProduceMessage(orderRequest.TransactionID, message)

	if errProduceMessage != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, errProduceMessage.Error()), errProduceMessage
	}

	return createOrderResponse(orderRequest, http.StatusCreated, "CREATED"), nil
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

func createOrderResponse(orderRequest domain.OrderRequest, httpStatusCode int, responseMessage string) domain.OrderResponse {
	orderResponse := domain.OrderResponse {
		OrderType: orderRequest.OrderType,
		OrderService: "orchestration",
		TransactionID: orderRequest.TransactionID,
		UserID: orderRequest.UserID,
		Action: "CREATE ORDER",
		ResponseCode: httpStatusCode,
		ResponseStatus: http.StatusText(httpStatusCode),
		ResponseMessage: responseMessage,
		Payload: orderRequest.Payload,
		ResponseCreatedAt: time.Now().Format("02-Jan-2006 15:04:05"),
	}
	return orderResponse
}

package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"order_service/internal/domain"
	"order_service/internal/domain/payload/request"
	"order_service/internal/provider/kafka"
	"order_service/internal/repository"
	"order_service/internal/usecase"
	"order_service/internal/utils"
	"os"
	"time"

	"github.com/rs/zerolog/log"
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
	log.Trace().Msg("Inside order usecase")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	log.Trace().Msg("Fetching last transactionID")
	lastID, err := uc.GetLastTransactionID()

	if err != nil {
		if err == sql.ErrNoRows {
			log.Trace().Msg("Set last ID if db is empty")
			lastID = ""
		} else {
			return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
		}
	}
	log.Trace().Msg("Trying to get generated ID")
	generatedID, err := utils.GenerateNextTransactionID(lastID)
	log.Debug().Msgf("Geneted ID is: %s", generatedID)

	if err != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
	}

	log.Trace().Msg("Setting generatedID as TransactionID")
	orderRequest.TransactionID = generatedID

	log.Trace().Msg("Creating transaction")
	transaction := domain.Transaction {
		TransactionID: generatedID,
		OrderType: orderRequest.OrderType,
		Status: "CREATED",
		UserID: orderRequest.UserID,
	}

	var transactionDetail domain.TransactionDetail

	log.Trace().Msg("Creating transaction details")
	if orderRequest.OrderType == "membership" {
		var membershipRequest request.MembershipRequest

		payloadBytes, err := json.Marshal(orderRequest.Payload)
		if err != nil {
			log.Error().Msgf("Failed to marshal payload: %v", err)
			return domain.OrderResponse{}, utils.ErrFailedToMarshallPayload
		}

		err = json.Unmarshal(payloadBytes, &membershipRequest)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal payload into MembershipRequest: %v", err)
			return domain.OrderResponse{}, utils.ErrFailedToUnmarshallPayload
		}

		transactionDetail = domain.TransactionDetail{
			TransactionID: generatedID,
			OrderType:     orderRequest.OrderType,
			Topic:         "orchestration",
			Action:        "START",
			Service:       "orchestration",
			Status:        "ONPROGRESS",
			Payload:       membershipRequest,
			UserID:        orderRequest.UserID,
			StatusCode: http.StatusCreated,
			StatusDesc: http.StatusText(http.StatusCreated),
			Message: "CREATED",
		}
	} else if orderRequest.OrderType == "donation" {
		var donationRequest request.DonationRequest

		payloadBytes, err := json.Marshal(orderRequest.Payload)
		if err != nil {
			log.Error().Msgf("Failed to marshal payload: %v", err)
			return domain.OrderResponse{}, utils.ErrFailedToMarshallPayload
		}

		err = json.Unmarshal(payloadBytes, &donationRequest)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal payload into MembershipRequest: %v", err)
			return domain.OrderResponse{}, utils.ErrFailedToUnmarshallPayload
		}

		transactionDetail = domain.TransactionDetail{
			TransactionID: generatedID,
			OrderType:     orderRequest.OrderType,
			Topic:         "orchestration",
			Action:        "START",
			Service:       "orchestration",
			Status:        "ONPROGRESS",
			Payload:       donationRequest,
			UserID:        orderRequest.UserID,
			StatusCode: http.StatusCreated,
			StatusDesc: http.StatusText(http.StatusCreated),
			Message: "CREATED",
		}
	}

	log.Trace().Msg("Calling order repo to save transaction")
	errSave := uc.repo.SaveTransaction(ctx, transaction)

	if errSave != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, errSave.Error()), errSave
	}

	log.Trace().Msg("Calling order repo to save transaction detail")
	errSaveDetail := uc.SaveTransactionDetails(transactionDetail)

	if errSaveDetail != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, fmt.Sprintf("Transaction created but couldn't write the details with message: %s", errSaveDetail.Error())), errSaveDetail
	}

	log.Trace().Msg("Trying to create new kafka producer instance")
	log.Debug().Msgf("Transaction detail topic is: %s", transactionDetail.Topic)
	broker := []string{os.Getenv("KAFKA_BROKER")}
	producer, err := kafka.NewKafkaProducer(broker, transactionDetail.Topic, 1, 1)

	if err != nil {
		return createOrderResponse(orderRequest, http.StatusInternalServerError, err.Error()), err
	}

	var message any

	if transactionDetail.OrderType == "membership" {
		// First, marshal the payload to a JSON byte slice
		payloadBytes, err := json.Marshal(orderRequest.Payload)
		if err != nil {
			log.Error().Msgf("Failed to marshal payload: %v", err)
			return createOrderResponse(orderRequest, http.StatusInternalServerError, "Failed to marshal payload"), err
		}

		// Create a MembershipRequest instance
		var membershipRequest request.MembershipRequest
		err = json.Unmarshal(payloadBytes, &membershipRequest)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal payload into MembershipRequest: %v", err)
			return createOrderResponse(orderRequest, http.StatusInternalServerError, "Failed to unmarshal payload"), err
		}

		// Create the message
		message = domain.Message[domain.TransactionMessage, request.MembershipRequest]{
			Header: domain.TransactionMessage{
				TransactionID: transactionDetail.TransactionID,
				OrderType:     transactionDetail.OrderType,
				UserID:        transactionDetail.UserID,
				Topic:         transactionDetail.Topic,
				Action:        transactionDetail.Action,
				Service:       transactionDetail.Service,
				Status:        transactionDetail.Status,
				StatusCode: transactionDetail.StatusCode,
				StatusDesc: transactionDetail.StatusDesc,
				Message: transactionDetail.Message,
			},
			Body: membershipRequest, // Use the deserialized membershipRequest
		}
	} else if transactionDetail.OrderType == "donation" {
		payloadBytes, err := json.Marshal(orderRequest.Payload)
		if err != nil {
			log.Error().Msgf("Failed to marshal payload: %v", err)
			return createOrderResponse(orderRequest, http.StatusInternalServerError, "Failed to marshal payload"), err
		}

		// Create a MembershipRequest instance
		var donationRequest request.DonationRequest
		err = json.Unmarshal(payloadBytes, &donationRequest)
		if err != nil {
			log.Error().Msgf("Failed to unmarshal payload into MembershipRequest: %v", err)
			return createOrderResponse(orderRequest, http.StatusInternalServerError, "Failed to unmarshal payload"), err
		}

		// Create the message
		message = domain.Message[domain.TransactionMessage, request.DonationRequest]{
			Header: domain.TransactionMessage{
				TransactionID: transactionDetail.TransactionID,
				OrderType:     transactionDetail.OrderType,
				UserID:        transactionDetail.UserID,
				Topic:         transactionDetail.Topic,
				Action:        transactionDetail.Action,
				Service:       transactionDetail.Service,
				Status:        transactionDetail.Status,
				StatusCode: transactionDetail.StatusCode,
				StatusDesc: transactionDetail.StatusDesc,
				Message: transactionDetail.Message,
			},
			Body: donationRequest, // Use the deserialized membershipRequest
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
	log.Trace().Msg("Calling get last transactionID repo")
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
	log.Trace().Msgf("ResponseCreatedAt: %s", orderResponse.ResponseCreatedAt)
	return orderResponse
}

package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/domain/payload/request"
	"orchestration_service/internal/provider/kafka"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/usecase"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type TransactionUseCase struct {
	repo repository.TransactionRepositoryInterface
}

func NewTransactionUseCase(repo repository.TransactionRepositoryInterface) usecase.TransactionUseCaseInterface {
	return TransactionUseCase {
		repo: repo,
	}
}

func (uc TransactionUseCase) TransactionUpdate(status string, transactionID string) error {
	log.Trace().Msg("Inside transaction usecase transaction update")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.TransactionUpdate(ctx, status, transactionID)
}

func (uc TransactionUseCase) TransactionDetailInput(transactionMessage domain.TransactionMessage, payload interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.TransactionDetailInput(ctx, transactionMessage, payload)
}

func (uc TransactionUseCase) FindTransactionDetailByIDStatusFailed(transactionID string) (domain.TransactionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.FindTransactionDetailByIDStatusFailed(ctx, transactionID)
}

func (uc TransactionUseCase) TransactionDetailRetry(id int, transactionDetail domain.TransactionDetail) (domain.TransactionDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return uc.repo.TransactionDetailRetry(ctx, id, transactionDetail)
}

func (uc TransactionUseCase) TransactionDetailSend(transactionDetail domain.TransactionDetail) {
	var messageSend any

	log.Trace().Msg("checking if order type is membership")
	switch transactionDetail.OrderType {
	case "membership":

		var body request.MembershipRequest
		payload, ok := transactionDetail.Payload.([]byte)
		if !ok {
			log.Error().Msg("Error: Payload is not of type []byte")
			return
		}

		// Proceed to unmarshal
		log.Trace().Msg("Trying to unmarshal payload to body")
		if err := json.Unmarshal(payload, &body); err != nil {
			log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
			return
		}

		log.Trace().Msg("Init message send")
		messageSend = domain.Message[domain.TransactionMessage, request.MembershipRequest] {
			Header: domain.TransactionMessage{
				TransactionID: transactionDetail.TransactionID,
				OrderType: transactionDetail.OrderType,
				UserID: transactionDetail.UserID,
				Topic: transactionDetail.Topic,
				Action: transactionDetail.Action,
				Service: transactionDetail.Service,
				Status: transactionDetail.Status,
				StatusCode: transactionDetail.StatusCode,
				StatusDesc: transactionDetail.StatusDesc,
				Message: transactionDetail.Message,
				CreatedAt: time.Now(),
			},
			Body: body,
		}

		log.Trace().Msg("Trying to create kafka producer")
		broker := []string{os.Getenv("KAFKA_BROKER")}
		producer, err := kafka.NewKafkaProducer(broker, transactionDetail.Topic, 1, 1)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
			return
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionDetail.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	case "donation":
		var body request.DonationRequest

		payload, ok := transactionDetail.Payload.([]byte)
		if !ok {
			log.Error().Msg("Error: Payload is not of type []byte")
			return
		}

		// Proceed to unmarshal
		log.Trace().Msg("Trying to unmarshal payload to body")
		if err := json.Unmarshal(payload, &body); err != nil {
			log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
			return
		}

		log.Trace().Msg("Init message send")
		messageSend = domain.Message[domain.TransactionMessage, request.DonationRequest]{
			Header: domain.TransactionMessage{
				TransactionID: transactionDetail.TransactionID,
				OrderType: transactionDetail.OrderType,
				UserID: transactionDetail.UserID,
				Topic: transactionDetail.Topic,
				Action: transactionDetail.Action,
				Service: transactionDetail.Service,
				Status: transactionDetail.Status,
				StatusCode: transactionDetail.StatusCode,
				StatusDesc: transactionDetail.StatusDesc,
				Message: transactionDetail.Message,
				CreatedAt: time.Now(),
			},
			Body: body,
		}

		log.Trace().Msg("Trying to create kafka producer")
		broker := []string{os.Getenv("KAFKA_BROKER")}
		producer, err := kafka.NewKafkaProducer(broker, transactionDetail.Topic, 1, 1)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
			return
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionDetail.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	}
}

func (uc TransactionUseCase) UpdateStatus(status string, transactionID string) error {
	log.Trace().Msg("Inside transaction usecase transaction UpdateStatus")
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	log.Trace().Msg("AAAAAAAAAAAAAAAAAAAAAAAAAAA")
	return uc.repo.UpdateStatus(ctx, status, transactionID)
}


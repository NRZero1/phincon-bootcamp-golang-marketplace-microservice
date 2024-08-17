package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/domain/payload/request"
	"orchestration_service/internal/provider/kafka"
	"orchestration_service/internal/repository"
	"orchestration_service/internal/usecase"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type ConsumerUseCase struct {
	configRepo repository.ConfigRepositoryInterface
}

func NewConsumerUseCase(configRepo repository.ConfigRepositoryInterface) usecase.ConsumerUseCaseInterface {
	return ConsumerUseCase{
		configRepo: configRepo,
	}
}

func (uc ConsumerUseCase) RouteMessage(message []byte) {
	var tempMessage struct {
		Header domain.TransactionMessage `json:"header"`
		Body   json.RawMessage           `json:"body"` // Keep the body as raw JSON to unmarshal later
	}

	if err := json.Unmarshal(message, &tempMessage); err != nil {
		log.Error().Msg(fmt.Sprintf("Error unmarshaling received message with error message: %s", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	transactionMessage := tempMessage.Header

	var statusCategory string
	if transactionMessage.StatusCode == http.StatusNotFound {
		statusCategory = "FAILED"
		// TODO update transaction table
	} else {
		statusCategory = "SUCCESS"
		// TODO update transaction table
	}

	configResponse, err := uc.configRepo.GetConfigByOrderType(ctx, transactionMessage.OrderType, transactionMessage.Service, statusCategory)

	if err != nil {
		return
	}

	var messageSend any

	switch transactionMessage.OrderType {
	case "membership":

		var body request.MembershipRequest
		if err := json.Unmarshal(tempMessage.Body, &body); err != nil {
			log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
			return
		}

		messageSend = domain.Message[domain.TransactionMessage, request.MembershipRequest] {
			Header: domain.TransactionMessage{
				TransactionID: transactionMessage.TransactionID,
				OrderType: transactionMessage.OrderType,
				UserID: transactionMessage.UserID,
				Topic: configResponse.ServiceDest,
				Action: configResponse.Action,
				Service: transactionMessage.Service,
				Status: transactionMessage.Status,
				StatusCode: transactionMessage.StatusCode,
				StatusDesc: transactionMessage.StatusDesc,
				Message: transactionMessage.Message,
				CreatedAt: time.Now(),
			},
			Body: body,
		}

		broker := []string{os.Getenv("KAFKA_BROKER")}
		producer, err := kafka.NewKafkaProducer(broker, configResponse.ServiceDest, 1, 1)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
			return
		}

		errProduceMessage := producer.ProduceMessage(transactionMessage.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	}
}

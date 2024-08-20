package impl

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"user_kafka_service/internal/domain"
	"user_kafka_service/internal/domain/payload/request"
	"user_kafka_service/internal/provider/kafka"
	"user_kafka_service/internal/usecase"
	"user_kafka_service/internal/utils"

	"github.com/rs/zerolog/log"
)

type ConsumerUseCase struct {
	userRequestUseCase usecase.UserRequestUseCaseInterface
}

func NewConsumerUseCase(userRequestUseCase usecase.UserRequestUseCaseInterface) usecase.ConsumerUseCaseInterface {
	return ConsumerUseCase{
		userRequestUseCase: userRequestUseCase,
	}
}

func (uc ConsumerUseCase) RouteMessage(message []byte) {
	log.Trace().Msg("Entering consumer use case route message")
	var tempMessage struct {
		Header domain.TransactionMessage `json:"header"`
		Body   json.RawMessage           `json:"body"` // Keep the body as raw JSON to unmarshal later
	}

	log.Trace().Msg("Trying to unmarshall message")
	if err := json.Unmarshal(message, &tempMessage); err != nil {
		log.Error().Msg(fmt.Sprintf("Error unmarshaling received message with error message: %s", err.Error()))
		return
	}

	transactionMessage := tempMessage.Header

	log.Trace().Msg("Trying to create kafka producer")
	broker := []string{os.Getenv("KAFKA_BROKER")}
	producer, err := kafka.NewKafkaProducer(broker, "orchestration", 1, 1)

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
		return
	}

	log.Trace().Msg("Trying to check message topic")
	if transactionMessage.Topic == "user-source" && transactionMessage.OrderType == "membership" {
		log.Trace().Msg("inside if topic is user-source and order type is membership")
		var membership request.MembershipRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &membership); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch user")
		_, statusCode, err := uc.userRequestUseCase.GetUserByID(transactionMessage.UserID)

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "user-source", "FAILED", statusCode, "request returned an error with not auto retryable status code", membership, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "user-source", "FAILED", statusCode, "request returned an error and system tried to retry but max retry reached", membership, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		sendMessage(transactionMessage, "user-source", "SUCCESS", statusCode, "SUCCESS", membership, producer)
		return
	} else if transactionMessage.Topic == "user-dest" && transactionMessage.OrderType == "membership" {
		log.Trace().Msg("inside if topic is user-dest and order type is membership")
		var membership request.MembershipRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &membership); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch user")
		_, statusCode, err := uc.userRequestUseCase.GetUserByID(int(membership.UserIDDest))

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "user-dest", "FAILED", statusCode, "Request returned an error with not auto retryable status code", membership, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "user-dest", "FAILED", statusCode, "Request returned an error and system tried to retry but max retry reached", membership, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		sendMessage(transactionMessage, "user-dest", "SUCCESS", statusCode, "SUCCESS", membership, producer)
		return
	} else if transactionMessage.Topic == "user-source" && transactionMessage.OrderType == "donation" {
		log.Trace().Msg("inside if topic is user-source and order type is donation")
		var donation request.DonationRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &donation); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch user")
		_, statusCode, err := uc.userRequestUseCase.GetUserByID(int(transactionMessage.UserID))

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "user-source", "FAILED", statusCode, "Request returned an error with not auto retryable status code", donation, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "user-source", "FAILED", statusCode, "Request returned an error and system tried to retry but max retry reached", donation, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		sendMessage(transactionMessage, "user-source", "SUCCESS", statusCode, "SUCCESS", donation, producer)
		return
	} else if transactionMessage.Topic == "user-dest" && transactionMessage.OrderType == "donation" {
		log.Trace().Msg("inside if topic is user-dest and order type is donation")
		var membership request.DonationRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &membership); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch user")
		_, statusCode, err := uc.userRequestUseCase.GetUserByID(int(membership.UserIDDest))

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "user-dest", "FAILED", statusCode, "Request returned an error with not auto retryable status code", membership, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "user-dest", "FAILED", statusCode, "Request returned an error and system tried to retry but max retry reached", membership, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		sendMessage(transactionMessage, "user-dest", "SUCCESS", statusCode, "SUCCESS", membership, producer)
		return
	}
}

func sendMessage(transactionMessage domain.TransactionMessage, service string, status string, statusCode int, message string, payload any, producer *kafka.KafkaProducer) {
	log.Trace().Msg("Trying to send message")
	transactionMessageSend := domain.TransactionMessage{
		TransactionID: transactionMessage.TransactionID,
		OrderType:     transactionMessage.OrderType,
		UserID:        transactionMessage.UserID,
		Topic:         "orchestration",
		Action:        transactionMessage.Action,
		Service:       service,
		Status:        status,
		StatusCode:    statusCode,
		StatusDesc:    http.StatusText(statusCode),
		Message:       message,
		CreatedAt:     time.Now(),
	}

var messageSend any

	if transactionMessage.OrderType == "membership" {
		log.Trace().Msg("If order type is membership")
		// Check the type assertion to prevent panics
		membershipRequest, ok := payload.(request.MembershipRequest)
		if !ok {
			log.Error().Msg("Payload is not of type request.MembershipRequest")
			return
		}

		messageSend = domain.Message[domain.TransactionMessage, request.MembershipRequest]{
			Header: transactionMessageSend,
			Body:   membershipRequest,
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionMessageSend.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	} else if transactionMessage.OrderType == "donation" {
		log.Trace().Msg("If order type is donation")
		// Check the type assertion to prevent panics
		donationRequest, ok := payload.(request.DonationRequest)
		if !ok {
			log.Error().Msg("Payload is not of type request.DonationRequest")
			return
		}

		messageSend = domain.Message[domain.TransactionMessage, request.DonationRequest]{
			Header: transactionMessageSend,
			Body:   donationRequest,
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionMessageSend.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	}

}

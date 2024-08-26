package impl

import (
	"channel_kafka_service/internal/domain"
	"channel_kafka_service/internal/domain/payload/request"
	"channel_kafka_service/internal/provider/kafka"
	"channel_kafka_service/internal/usecase"
	"channel_kafka_service/internal/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type ConsumerUseCase struct {
	channelRequestUseCase usecase.ChannelRequestUseCaseInterface
}

func NewConsumerUseCase(channelRequestUseCase usecase.ChannelRequestUseCaseInterface) usecase.ConsumerUseCaseInterface {
	return ConsumerUseCase{
		channelRequestUseCase: channelRequestUseCase,
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
	if transactionMessage.Topic == "channel" && transactionMessage.OrderType == "membership" {
		log.Trace().Msg("inside if topic is channel and order type is membership")
		var membership request.MembershipRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &membership); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch channel")
		_, statusCode, resp, err := uc.channelRequestUseCase.GetChannelByID(int(membership.ChannelID))

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "FAILED", statusCode, "Request returned an error with not auto retryable status code", membership, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "FAILED", statusCode, "Request returned an error and system tried to retry but max retry reached", membership, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		membership.UserIDDest = resp.UserID
		log.Debug().Msgf("Membership value: %+v", membership)
		sendMessage(transactionMessage, "SUCCESS", statusCode, "SUCCESS", membership, producer)
		return
	} else if transactionMessage.Topic == "channel" && transactionMessage.OrderType == "donation" {
		log.Trace().Msg("inside if topic is channel and order type is donation")
		var donation request.DonationRequest

		log.Trace().Msg("Trying to unmarshal message body")
		if err := json.Unmarshal(tempMessage.Body, &donation); err != nil {
            log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
            return
        }

		log.Trace().Msg("Trying to fetch channel")
		_, statusCode, resp, err := uc.channelRequestUseCase.GetChannelByID(int(donation.ChannelID))

		log.Trace().Msg("Trying to check if error is exist")
		if err != nil {
			if errors.Is(err, utils.ErrHttpNotRetryAble) {
				log.Trace().Msg("Error is not retryable")
				sendMessage(transactionMessage, "FAILED", statusCode, "Request returned an error with not auto retryable status code", donation, producer)
			} else if errors.Is(err, utils.ErrMaxRetryReached) {
				log.Trace().Msg("Max retry reached")
				sendMessage(transactionMessage, "FAILED", statusCode, "Request returned an error and system tried to retry but max retry reached", donation, producer)
			}
			return
		}
		log.Trace().Msg("Fetch success, returning")
		donation.UserIDDest = resp.UserID
		log.Debug().Msgf("Donation value: %+v", donation)
		sendMessage(transactionMessage, "SUCCESS", statusCode, "SUCCESS", donation, producer)
		return
	}
}

func sendMessage(transactionMessage domain.TransactionMessage, status string, statusCode int, message string, payload any, producer *kafka.KafkaProducer) {
	log.Trace().Msg("Trying to send message")
	transactionMessageSend := domain.TransactionMessage{
		TransactionID: transactionMessage.TransactionID,
		OrderType:     transactionMessage.OrderType,
		UserID:        transactionMessage.UserID,
		Topic:         "orchestration",
		Action:        transactionMessage.Action,
		Service:       "channel",
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

		// Serialize messageSend to JSON
		// serializedMessage, err := json.Marshal(messageSend)
		// if err != nil {
		// 	log.Error().Msg(fmt.Sprintf("Error marshalling messageSend: %s", err.Error()))
		// 	return
		// }

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

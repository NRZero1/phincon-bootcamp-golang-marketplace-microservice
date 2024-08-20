package impl

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/domain/payload/request"
	"orchestration_service/internal/provider/kafka"
	"orchestration_service/internal/usecase"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

type ConsumerUseCase struct {
	configUseCase usecase.ConfigUseCaseInterface
	transactionUseCase usecase.TransactionUseCaseInterface
}

func NewConsumerUseCase(configUseCase usecase.ConfigUseCaseInterface, transactionUseCase usecase.TransactionUseCaseInterface) usecase.ConsumerUseCaseInterface {
	return ConsumerUseCase{
		configUseCase: configUseCase,
		transactionUseCase: transactionUseCase,
	}
}

func (uc ConsumerUseCase) RouteMessage(message []byte) {
	log.Trace().Msg("Inside Consumer Use Case Route message")
	var tempMessage struct {
		Header domain.TransactionMessage `json:"header"`
		Body   json.RawMessage           `json:"body"` // Keep the body as raw JSON to unmarshal later
	}

	log.Trace().Msg("Trying to unmarshal message received")
	if err := json.Unmarshal(message, &tempMessage); err != nil {
		log.Error().Msg(fmt.Sprintf("Error unmarshaling received message with error message: %s", err.Error()))
		return
	}

	transactionMessage := tempMessage.Header

	var statusCategory string
	log.Trace().Msg("Trying to check status code and error if exist")
	log.Debug().Msgf("Transaction status is %s", transactionMessage.Status)
	if transactionMessage.Status == "FAILED" {
		log.Trace().Msg("Checking order type when failed")
		if transactionMessage.OrderType == "membership" {
			log.Trace().Msg("If order type is membership")
			var body request.MembershipRequest

			log.Trace().Msg("Trying to unmarshal temp body to body")
			if err := json.Unmarshal(tempMessage.Body, &body); err != nil {
				log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
				return
			}

			log.Debug().Msgf("Is payment pending: %t", body.IsPaymentPending)
			if body.IsPaymentPending {
				log.Trace().Msg("Check if payment is pending")
				messageSend := domain.Message[domain.TransactionMessage, request.MembershipRequest] {
					Header: domain.TransactionMessage{
						TransactionID: transactionMessage.TransactionID,
						OrderType: transactionMessage.OrderType,
						UserID: transactionMessage.UserID,
						Topic: "balance",
						Action: "ADD BALANCE",
						Service: "orchestration",
						Status: transactionMessage.Status,
						StatusCode: transactionMessage.StatusCode,
						StatusDesc: transactionMessage.StatusDesc,
						Message: transactionMessage.Message,
						CreatedAt: time.Now(),
					},
					Body: body,
				}

				log.Trace().Msg("Trying to call transaction use case to update transaction")
				errTransactionUpdate := uc.transactionUseCase.TransactionUpdate(transactionMessage.Status, transactionMessage.TransactionID)

				if errTransactionUpdate != nil {
					log.Error().Msgf("failed to update transaction with error message: %s", errTransactionUpdate.Error())
					return
				}

				log.Trace().Msg("Trying to call transaction use case to insert transaction detail")
				errTransactionDetailInput := uc.transactionUseCase.TransactionDetailInput(transactionMessage, tempMessage.Body)

				if errTransactionDetailInput != nil {
					log.Error().Msgf("failed to update transaction with error message: %s", errTransactionDetailInput.Error())
					return
				}

				broker := []string{os.Getenv("KAFKA_BROKER")}
				producer, err := kafka.NewKafkaProducer(broker, "balance", 1, 1)

				if err != nil {
					log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
					return
				}

				errProduceMessage := producer.ProduceMessage(transactionMessage.TransactionID, messageSend)

				if errProduceMessage != nil {
					log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
					return
				}

				log.Printf("===============================================================================================")
				log.Printf("Transaction ID %s for order type '%s' is FAILED, User Balance is ROLLBACKED\n", transactionMessage.TransactionID, transactionMessage.OrderType)
				log.Printf("===============================================================================================")
				return
			}
		} else if transactionMessage.OrderType == "donation" {
			log.Trace().Msg("If order type is donation")
			var body request.DonationRequest

			log.Trace().Msg("Trying to unmarshal temp body to body")
			if err := json.Unmarshal(tempMessage.Body, &body); err != nil {
				log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
				return
			}

			log.Debug().Msgf("Is payment pending: %t", body.IsPaymentPending)
			if body.IsPaymentPending {
				log.Trace().Msg("Check if payment is pending")
				messageSend := domain.Message[domain.TransactionMessage, request.DonationRequest]{
					Header: domain.TransactionMessage{
						TransactionID: transactionMessage.TransactionID,
						OrderType: transactionMessage.OrderType,
						UserID: transactionMessage.UserID,
						Topic: "balance",
						Action: "ADD BALANCE",
						Service: "orchestration",
						Status: transactionMessage.Status,
						StatusCode: transactionMessage.StatusCode,
						StatusDesc: transactionMessage.StatusDesc,
						Message: transactionMessage.Message,
						CreatedAt: time.Now(),
					},
					Body: body,
				}

				log.Trace().Msg("Trying to call transaction use case to update transaction")
				errTransactionUpdate := uc.transactionUseCase.TransactionUpdate(transactionMessage.Status, transactionMessage.TransactionID)

				if errTransactionUpdate != nil {
					log.Error().Msgf("failed to update transaction with error message: %s", errTransactionUpdate.Error())
					return
				}

				log.Trace().Msg("Trying to call transaction use case to insert transaction detail")
				errTransactionDetailInput := uc.transactionUseCase.TransactionDetailInput(transactionMessage, tempMessage.Body)

				if errTransactionDetailInput != nil {
					log.Error().Msgf("failed to update transaction with error message: %s", errTransactionDetailInput.Error())
					return
				}

				broker := []string{os.Getenv("KAFKA_BROKER")}
				producer, err := kafka.NewKafkaProducer(broker, "balance", 1, 1)

				if err != nil {
					log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
					return
				}

				errProduceMessage := producer.ProduceMessage(transactionMessage.TransactionID, messageSend)

				if errProduceMessage != nil {
					log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
					return
				}

				log.Printf("===============================================================================================")
				log.Printf("Transaction ID %s for order type '%s' is FAILED, User Balance is ROLLBACKED\n", transactionMessage.TransactionID, transactionMessage.OrderType)
				log.Printf("===============================================================================================")
				return
			}
		}
	} else {
		statusCategory = "SUCCESS"
		log.Trace().Msg("Checking order type when success")
		if transactionMessage.OrderType == "membership" {
			log.Trace().Msg("Trying to call transaction use case to update transaction")
			errTransactionUpdate := uc.transactionUseCase.TransactionUpdate("SUCCESS", transactionMessage.TransactionID)

			if errTransactionUpdate != nil {
				log.Error().Msgf("failed to update transaction with error message: %s", errTransactionUpdate.Error())
				return
			}
			if transactionMessage.Service == "balance-dest" && transactionMessage.Action == "ADD BALANCE" {
				log.Printf("===============================================================================================")
				log.Printf("Transaction ID %s for order type '%s' is COMPLETED\n", transactionMessage.TransactionID, transactionMessage.OrderType)
				log.Printf("===============================================================================================")
				return
			}
		} else if transactionMessage.OrderType == "donation" {
			if transactionMessage.Service == "balance-dest" && transactionMessage.Action == "ADD BALANCE" {
				log.Printf("===============================================================================================")
				log.Printf("Transaction ID %s for order type '%s' is COMPLETED\n", transactionMessage.TransactionID, transactionMessage.OrderType)
				log.Printf("===============================================================================================")
				return
			}
		}
	}

	log.Trace().Msg("Trying to call transaction use case to update transaction")
	errTransactionUpdate := uc.transactionUseCase.TransactionUpdate("ONPROGRESS", transactionMessage.TransactionID)

	if errTransactionUpdate != nil {
		log.Error().Msgf("failed to update transaction with error message: %s", errTransactionUpdate.Error())
		return
	}

	log.Trace().Msg("Trying to call transaction use case to insert transaction detail")
	errTransactionDetailInput := uc.transactionUseCase.TransactionDetailInput(transactionMessage, tempMessage.Body)

	if errTransactionDetailInput != nil {
		log.Error().Msgf("failed to update transaction with error message: %s", errTransactionDetailInput.Error())
		return
	}

	log.Trace().Msg("Trying to call config use case to fetch config")
	configResponse, err := uc.configUseCase.GetConfigByOrderType(transactionMessage.OrderType, transactionMessage.Service, statusCategory)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return
		} else if errors.Is(err, sql.ErrNoRows) && transactionMessage.Status == "SUCCESS" {
			log.Info().Msgf("Transaction with transactionID %s completed", transactionMessage.TransactionID)
			return
		}
		log.Error().Msg(fmt.Sprintf("Error fetching config: %s", err.Error()))
		return
	}

	var messageSend any

	log.Trace().Msg("checking if order type is membership")
	switch transactionMessage.OrderType {
	case "membership":

		var body request.MembershipRequest
		log.Trace().Msg("Trying to unmarshal temp body to body")
		if err := json.Unmarshal(tempMessage.Body, &body); err != nil {
			log.Error().Msg(fmt.Sprintf("Error unmarshaling membership request body: %s", err.Error()))
			return
		}

		log.Trace().Msg("Init message send")
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

		log.Trace().Msg("Trying to create kafka producer")
		broker := []string{os.Getenv("KAFKA_BROKER")}
		producer, err := kafka.NewKafkaProducer(broker, configResponse.ServiceDest, 1, 1)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
			return
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionMessage.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	case "donation":
		var body request.DonationRequest
		log.Trace().Msg("Trying to unmarshal temp body to body")
		if err := json.Unmarshal(tempMessage.Body, &body); err != nil {
			log.Error().Msg(fmt.Sprintf("Error unmarshaling donation request body: %s", err.Error()))
			return
		}

		log.Trace().Msg("Init message send")
		messageSend = domain.Message[domain.TransactionMessage, request.DonationRequest]{
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

		log.Trace().Msg("Trying to create kafka producer")
		broker := []string{os.Getenv("KAFKA_BROKER")}
		producer, err := kafka.NewKafkaProducer(broker, configResponse.ServiceDest, 1, 1)

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create new kafka producer with error message: %s", err.Error()))
			return
		}

		log.Trace().Msg("Trying to send message")
		errProduceMessage := producer.ProduceMessage(transactionMessage.TransactionID, messageSend)

		if errProduceMessage != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to send message with error message: %s", errProduceMessage.Error()))
			return
		}
	}
}

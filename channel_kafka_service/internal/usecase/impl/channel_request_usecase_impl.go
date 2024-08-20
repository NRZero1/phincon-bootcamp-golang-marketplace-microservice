package impl

import (
	response "channel_kafka_service/internal/domain/channel_service_response"
	"channel_kafka_service/internal/usecase"
	"channel_kafka_service/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type ChannelRequestUseCase struct {
}

func NewChannelRequestUseCase() usecase.ChannelRequestUseCaseInterface {
	return ChannelRequestUseCase{}
}

func (uc ChannelRequestUseCase) GetChannelByID(id int) (bool, int, response.Channel, error) {
	log.Trace().Msg("Inside channel request use case")
	counter := 1
	retryAfter := 5
	var globalResponse response.GlobalResponse
	for {
		log.Trace().Msg("Inside infinite loop to retry request")
		log.Debug().Msgf("Counter value is %d", counter)
		if counter == 4 {
			return false, globalResponse.StatusCode, response.Channel{}, utils.ErrMaxRetryReached
		}

		log.Trace().Msg("Trying to send request")
		resp, err := http.Get(fmt.Sprintf("http://localhost:8084/channel/%d", id))

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create request to channel service with message: %s", err.Error()))
			return false, http.StatusInternalServerError, response.Channel{}, utils.ErrHttpRequest
		}

		defer resp.Body.Close()

		log.Trace().Msg("Trying to decode json")
		errDecode := json.NewDecoder(resp.Body).Decode(&globalResponse)
		log.Debug().Msgf("global response value: %+v", globalResponse)

		if errDecode != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to decode request channel response with error message: %s", errDecode.Error()))
			return false, http.StatusInternalServerError, response.Channel{}, utils.ErrJsonDecode
		}

		if globalResponse.Data != nil {
			if utils.IsRetryAbleStatusCode(globalResponse.StatusCode) {
				log.Info().Msgf("Request result an error, retrying in %d seconds", retryAfter * counter)
				time.Sleep(time.Duration(retryAfter * counter) * time.Second)
				counter++
				continue
			} else if utils.IsNotRetryAble(globalResponse.StatusCode) {
				log.Error().Msg("Request result an error and system cannot retry, need to edit and manual")
				return false, globalResponse.StatusCode, response.Channel{}, utils.ErrHttpNotRetryAble
			} else if utils.IsSuccess(globalResponse.StatusCode) {
				var body response.Channel

				// Handle different possible types of globalResponse.Data
				switch data := globalResponse.Data.(type) {
				case map[string]interface{}:
					log.Trace().Msg("Data is a map, trying to marshal it back to JSON")
					// Marshal the map back to JSON so it can be unmarshalled into the response.Channel struct
					jsonData, err := json.Marshal(data)
					if err != nil {
						log.Error().Msg(fmt.Sprintf("Error marshaling globalResponse.Data: %s", err.Error()))
						return false, http.StatusInternalServerError, response.Channel{}, err
					}
					log.Trace().Msg("Trying to unmarshal JSON data into body")
					if err := json.Unmarshal(jsonData, &body); err != nil {
						log.Error().Msg(fmt.Sprintf("Error unmarshaling JSON data into response.Channel: %s", err.Error()))
						return false, http.StatusInternalServerError, response.Channel{}, utils.ErrJsonDecode
					}

				case []byte:
					log.Trace().Msg("Data is raw JSON bytes, trying to unmarshal into body")
					if err := json.Unmarshal(data, &body); err != nil {
						log.Error().Msg(fmt.Sprintf("Error unmarshaling raw JSON bytes into response.Channel: %s", err.Error()))
						return false, http.StatusInternalServerError, response.Channel{}, utils.ErrJsonDecode
					}

				default:
					log.Error().Msg("Unhandled type for globalResponse.Data")
					return false, http.StatusInternalServerError, response.Channel{}, utils.ErrTypeAssertion
				}

				return true, globalResponse.StatusCode, body, nil
			}
		}
	}
}

// func (uc ChannelRequestUseCase) AddMembership(channelID int, userID int) (bool, int, error) {
// 	log.Trace().Msg("Inside channel request use case")
// 	counter := 1
// 	retryAfter := 5
// 	var globalResponse response.GlobalResponse
// 	for {
// 		log.Trace().Msg("Inside infinite loop to retry request")
// 		log.Debug().Msgf("Counter value is %d", counter)
// 		if counter == 4 {
// 			return false, globalResponse.StatusCode, utils.ErrMaxRetryReached
// 		}

// 		log.Trace().Msg("Trying to send request")
// 		resp, err := http.Get(fmt.Sprintf("http://localhost:8081/channel/%d/", id))

// 		if err != nil {
// 			log.Error().Msg(fmt.Sprintf("Error when trying to create request to channel service with message: %s", err.Error()))
// 			return false, http.StatusInternalServerError, utils.ErrHttpRequest
// 		}

// 		defer resp.Body.Close()

// 		log.Trace().Msg("Trying to decode json")
// 		errDecode := json.NewDecoder(resp.Body).Decode(&globalResponse)

// 		if errDecode != nil {
// 			log.Error().Msg(fmt.Sprintf("Error when trying to decode request channel response with error message: %s", errDecode.Error()))
// 			return false, http.StatusInternalServerError, utils.ErrJsonDecode
// 		}

// 		if globalResponse.Data != nil {
// 			if utils.IsRetryAbleStatusCode(globalResponse.StatusCode) {
// 				log.Info().Msgf("Request result an error, retrying in %d seconds", retryAfter * counter)
// 				time.Sleep(time.Duration(retryAfter * counter) * time.Second)
// 				counter++
// 				continue
// 			} else if utils.IsNotRetryAble(globalResponse.StatusCode) {
// 				log.Error().Msg("Request result an error and system cannot retry, need to edit and manual")
// 				return false, globalResponse.StatusCode, utils.ErrHttpNotRetryAble
// 			} else if utils.IsSuccess(globalResponse.StatusCode) {
// 				return true, globalResponse.StatusCode, nil
// 			}
// 		}
// 	}
// }


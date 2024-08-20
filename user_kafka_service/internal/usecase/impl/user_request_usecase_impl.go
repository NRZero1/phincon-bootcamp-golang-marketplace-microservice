package impl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	response "user_kafka_service/internal/domain/user_service_response"
	"user_kafka_service/internal/usecase"
	"user_kafka_service/internal/utils"

	"github.com/rs/zerolog/log"
)

type UserRequestUseCase struct {
}

func NewUserRequestUseCase() usecase.UserRequestUseCaseInterface {
	return UserRequestUseCase{}
}

func (uc UserRequestUseCase) GetUserByID(id int) (bool, int, error) {
	log.Trace().Msg("Inside user request use case")
	counter := 1
	retryAfter := 5
	var globalResponse response.GlobalResponse
	for {
		log.Trace().Msg("Inside infinite loop to retry request")
		log.Debug().Msgf("Counter value is %d", counter)
		if counter == 4 {
			return false, globalResponse.StatusCode, utils.ErrMaxRetryReached
		}

		log.Trace().Msg("Trying to send request")
		resp, err := http.Get(fmt.Sprintf("http://localhost:8080/user/%d", id))

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create request to user service with message: %s", err.Error()))
			return false, http.StatusInternalServerError, utils.ErrHttpRequest
		}

		defer resp.Body.Close()

		log.Trace().Msg("Trying to decode json")
		errDecode := json.NewDecoder(resp.Body).Decode(&globalResponse)

		if errDecode != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to decode request user response with error message: %s", errDecode.Error()))
			return false, http.StatusInternalServerError, utils.ErrJsonDecode
		}

		if globalResponse.Data != nil {
			if utils.IsRetryAbleStatusCode(globalResponse.StatusCode) {
				log.Info().Msgf("Request result an error, retrying in %d seconds", retryAfter * counter)
				time.Sleep(time.Duration(retryAfter * counter) * time.Second)
				counter++
				continue
			} else if utils.IsNotRetryAble(globalResponse.StatusCode) {
				log.Error().Msg("Request result an error and system cannot retry, need to edit and manual")
				return false, globalResponse.StatusCode, utils.ErrHttpNotRetryAble
			} else if utils.IsSuccess(globalResponse.StatusCode) {
				return true, globalResponse.StatusCode, nil
			}
		}
	}
}

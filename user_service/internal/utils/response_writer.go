package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"user_service/internal/domain/dto/response"

	"github.com/rs/zerolog/log"
)

type Response struct {
	response.GlobalResponse
}

func NewResponse(message string, statusCode int, data interface{}) *Response {
	return &Response{
		GlobalResponse: response.GlobalResponse{
			Message: message,
			StatusCode: statusCode,
			Data: data,
		},
	}
}

func (r Response) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.GlobalResponse.StatusCode)
	
	res, err := json.MarshalIndent(&r.GlobalResponse, "", "	")

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to marshal response with error: %s", err.Error()))
		response := Response {
			GlobalResponse: response.GlobalResponse{
				Message: fmt.Sprintf("Error when trying to marshal response with error: %s", err.Error()),
				StatusCode: http.StatusInternalServerError,
				Data: nil,
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Write(res)
}
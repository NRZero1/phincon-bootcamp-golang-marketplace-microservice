package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"order_service/internal/domain"

	"github.com/rs/zerolog/log"
)

type ResponseOrder struct {
	domain.OrderResponse
}

type Response struct {
	domain.GlobalResponse
}

func NewOrderResponse(orderResponse domain.OrderResponse) *ResponseOrder {
	return &ResponseOrder{
		OrderResponse: domain.OrderResponse{
			OrderType: orderResponse.OrderType,
			OrderService: orderResponse.OrderService,
			TransactionID: orderResponse.TransactionID,
			UserID: orderResponse.UserID,
			ResponseCode: orderResponse.ResponseCode,
			ResponseStatus: orderResponse.ResponseStatus,
			ResponseMessage: orderResponse.ResponseMessage,
			Action: orderResponse.Action,
			Payload: orderResponse.Payload,
			ResponseCreatedAt: orderResponse.ResponseCreatedAt,
		},
	}
}

func (r ResponseOrder) WriteOrder(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.OrderResponse.ResponseCode)

	res, err := json.MarshalIndent(&r.OrderResponse, "", "	")

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to marshal response with error: %s", err.Error()))
		response := ResponseOrder {
			OrderResponse: domain.OrderResponse{
				OrderType: r.OrderType,
				OrderService: r.OrderService,
				TransactionID: r.TransactionID,
				UserID: r.UserID,
				ResponseCode: http.StatusInternalServerError,
				ResponseStatus: http.StatusText(http.StatusInternalServerError),
				ResponseMessage: fmt.Sprintf("Error when trying to marshal response with error: %s", err.Error()),
			},
		}

		json.NewEncoder(w).Encode(response)
		return
	}

	w.Write(res)
}

func NewResponse(message string, statusCode int, data interface{}) *Response {
	return &Response{
		GlobalResponse: domain.GlobalResponse{
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
			GlobalResponse: domain.GlobalResponse{
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

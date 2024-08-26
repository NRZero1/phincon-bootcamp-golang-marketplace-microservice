package impl

import (
	"net/http"
	"order_service/internal/domain"
	"order_service/internal/handler"
	"order_service/internal/usecase"
	"order_service/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log" // Import zerolog
)

type OrderHandler struct {
	usecase usecase.OrderUseCaseInterface
}

func NewOrderHandler(usecase usecase.OrderUseCaseInterface) handler.OrderHandlerInterface {
	return OrderHandler{
		usecase: usecase,
	}
}

func (h OrderHandler) Order(c *gin.Context) {
	log.Trace().Msg("Entering order handler")
	var orderRequest domain.OrderRequest

	// Log the incoming request
	log.Info().Msg("Received new order request")

	log.Trace().Msg("Trying to decode json")
	err := c.ShouldBindJSON(&orderRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON to orderRequest")

		orderResponse := domain.OrderResponse{
			OrderType:        orderRequest.OrderType,
			OrderService:     "order_service",
			UserID:           orderRequest.UserID,
			Action:           "CREATE ORDER",
			ResponseCode:     http.StatusBadRequest,
			ResponseStatus:   http.StatusText(http.StatusBadRequest),
			ResponseMessage:  utils.ErrJsonDecode.Error(),
			Payload:          orderRequest.Payload,
			ResponseCreatedAt: time.Now().Format("02-Jan-2006 15:04:05"),
		}

		log.Info().Msg("Sending bad request response due to JSON binding error")
		utils.NewOrderResponse(orderResponse).WriteOrder(c.Writer)
		return
	}

	log.Info().
		Str("OrderType", orderRequest.OrderType).
		Int("UserID", orderRequest.UserID).
		Msg("Processing order request")

	log.Trace().Msg("Calling order usecase")
		response, errSave := h.usecase.SaveTransaction(orderRequest)
	if errSave != nil {
		log.Error().Err(errSave).Msg("Failed to save transaction")
		utils.NewOrderResponse(response).WriteOrder(c.Writer)
		return
	}

	log.Info().Msg("Order transaction saved successfully")
	utils.NewOrderResponse(response).WriteOrder(c.Writer)
}

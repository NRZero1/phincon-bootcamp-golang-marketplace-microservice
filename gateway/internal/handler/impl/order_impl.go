package impl

import (
	"gateway/internal/domain"
	"gateway/internal/handler"
	"gateway/internal/usecase"
	"gateway/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	usecase usecase.OrderUseCaseInterface
}

func NewOrderHandler(usecase usecase.OrderUseCaseInterface) handler.OrderHandlerInterface {
	return OrderHandler {
		usecase: usecase,
	}
}

func (h OrderHandler) Order(c *gin.Context) {
	var orderRequest domain.OrderRequest
	err := c.ShouldBindJSON(&orderRequest)

	if err != nil {
		orderResponse := domain.OrderResponse {
			OrderType: orderRequest.OrderType,
			OrderService: "gateway",
			UserID: orderRequest.UserID,
			ResponseCode: http.StatusBadRequest,
			ResponseStatus: http.StatusText(http.StatusBadRequest),
			ResponseMessage: utils.ErrJsonDecode.Error(),
			Payload: orderRequest.Payload,
		}

		utils.NewOrderResponse(orderResponse).WriteOrder(c.Writer)
		return
	}

	errSave := h.usecase.SaveTransaction(orderRequest)

	if errSave != nil {
		orderResponse := domain.OrderResponse {
			OrderType: orderRequest.OrderType,
			OrderService: "gateway",
			UserID: orderRequest.UserID,
			ResponseCode: http.StatusInternalServerError,
			ResponseStatus: http.StatusText(http.StatusInternalServerError),
			ResponseMessage: errSave.Error(),
			Payload: orderRequest.Payload,
		}
		utils.NewOrderResponse(orderResponse).WriteOrder(c.Writer)
		return
	}

	orderResponse := domain.OrderResponse {
		OrderType: orderRequest.OrderType,
		OrderService: "gateway",
		UserID: orderRequest.UserID,
		ResponseCode: http.StatusCreated,
		ResponseStatus: http.StatusText(http.StatusCreated),
		ResponseMessage: "CREATED",
		Payload: orderRequest.Payload,
	}

	utils.NewOrderResponse(orderResponse).WriteOrder(c.Writer)
}
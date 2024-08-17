package impl

import (
	"gateway/internal/domain"
	"gateway/internal/handler"
	"gateway/internal/usecase"
	"gateway/internal/utils"
	"net/http"
	"time"

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
			Action: "CREATE ORDER",
			ResponseCode: http.StatusBadRequest,
			ResponseStatus: http.StatusText(http.StatusBadRequest),
			ResponseMessage: utils.ErrJsonDecode.Error(),
			Payload: orderRequest.Payload,
			ResponseCreatedAt: time.Now().Format("02-Jan-2006 15:04:05"),
		}

		utils.NewOrderResponse(orderResponse).WriteOrder(c.Writer)
		return
	}

	response, errSave := h.usecase.SaveTransaction(orderRequest)

	if errSave != nil {
		utils.NewOrderResponse(response).WriteOrder(c.Writer)
		return
	}

	utils.NewOrderResponse(response).WriteOrder(c.Writer)
}

func FindByTransactionID(context *gin.Context) {

}

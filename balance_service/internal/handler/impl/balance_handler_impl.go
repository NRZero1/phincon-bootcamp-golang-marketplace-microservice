package impl

import (
	"balance_service/internal/domain/dto/response"
	"balance_service/internal/handler"
	"balance_service/internal/usecase"
	"balance_service/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type BalanceHandler struct {
	usecase usecase.BalanceUseCaseInterface
}

func NewBalanceHandler(usecase usecase.BalanceUseCaseInterface) handler.BalanceHandlerInterface {
	return BalanceHandler{
		usecase: usecase,
	}
}

func (h BalanceHandler) FindByID(c *gin.Context) {
	log.Trace().Msg("Entering balance handler find by id")

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse{
			Message:    utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	foundBalance, errFound := h.usecase.FindByID(id)

	if errFound != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse

		log.Trace().Msg("Fetch error")
		log.Error().Str("Error message: ", errFound.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)

		resp = response.GlobalResponse{
			Message:    errFound.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       foundBalance,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h BalanceHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering balance get all handler")

	allBalances := h.usecase.GetAll()

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       allBalances,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}

// func (h BalanceHandler) GetAllOrFindByName(c *gin.Context) {
// 	name := c.Query("name")

// 	if name != "" {
// 		h.FindByName(c)
// 	} else {
// 		h.GetAll(c)
// 	}
// }

func (h BalanceHandler) Deduct(c *gin.Context) {
	log.Trace().Msg("Entering balance handler find by id")

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse{
			Message:    utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	amountString := c.Query("amount")

	amount, errConv := strconv.ParseFloat(amountString, 64)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting amount string to float")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse{
			Message:    utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	errDeduct := h.usecase.Deduct(id, amount)

	if errDeduct != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse

		log.Trace().Msg("Fetch error")
		log.Error().Str("Error message: ", errDeduct.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)

		resp = response.GlobalResponse{
			Message:    errDeduct.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       "OK",
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h BalanceHandler) AddBalance(c *gin.Context) {
	log.Trace().Msg("Entering balance handler find by id")

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse{
			Message:    utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	amountString := c.Query("amount")

	amount, errConv := strconv.ParseFloat(amountString, 64)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting amount string to float")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse{
			Message:    utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	errAdd := h.usecase.AddBalance(id, amount)

	if errAdd != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse

		log.Trace().Msg("Fetch error")
		log.Error().Str("Error message: ", errAdd.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)

		resp = response.GlobalResponse{
			Message:    errAdd.Error(),
			StatusCode: http.StatusBadRequest,
			Data:       "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       "OK",
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}


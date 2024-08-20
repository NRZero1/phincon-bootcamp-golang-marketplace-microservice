package impl

import (
	"errors"
	"fmt"
	"net/http"
	"orchestration_service/internal/domain"
	"orchestration_service/internal/handler"
	"orchestration_service/internal/usecase"
	"orchestration_service/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/rs/zerolog/log"
)

type TransactionHandler struct {
	usecase usecase.TransactionUseCaseInterface
}

func NewTransactionHandler(usecase usecase.TransactionUseCaseInterface) handler.TransactionHandlerInterface {
	return TransactionHandler {
		usecase: usecase,
	}
}

func (h TransactionHandler) FindTransactionDetailByIDStatusFailed(c *gin.Context) {
	log.Trace().Msg("Inside FindTransactionDetailByIDStatusFailed")
	transactionID := c.Param("transaction_id")
	log.Debug().Str("Received transactionID is: ", transactionID).Msg("debug")

	foundTransactionDetail, err := h.usecase.FindTransactionDetailByIDStatusFailed(transactionID)

	if err != nil {
		if errors.Is(err, utils.ErrNoSqlRows) {
			utils.NewResponse(err.Error(), http.StatusNotFound, nil).Write(c.Writer)
			return
		} else {
			utils.NewResponse(err.Error(), http.StatusInternalServerError, nil).Write(c.Writer)
			return
		}
	}

	utils.NewResponse("OK", http.StatusOK, foundTransactionDetail).Write(c.Writer)
}

func (h TransactionHandler) TransactionDetailRetry(c *gin.Context) {
	var transactionDetail domain.TransactionDetail

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		utils.NewResponse(http.StatusText(http.StatusBadRequest), http.StatusBadRequest, nil).Write(c.Writer)
		return
	}

	err := c.ShouldBindJSON(&transactionDetail)

	if err != nil {
		log.Error().Msgf("Failed to decode JSON with error message: %s", err.Error())
		utils.NewResponse(utils.ErrJsonDecode.Error(), http.StatusBadRequest, nil)
		return
	}

	log.Trace().Msg("Validating user input")
	errValidate := utils.ValidateStruct(&transactionDetail)

	if errValidate != nil {
		errors := make(map[string]string)
		if _, ok := errValidate.(*validator.InvalidValidationError); ok {
			log.Trace().Msg("Error with validator")
			log.Error().Str("Error message: ", errValidate.Error())
			c.JSON(http.StatusInternalServerError, errValidate.Error())
			return
		}

		for _, err := range errValidate.(validator.ValidationErrors) {
			errors[err.Field()] = fmt.Sprintf("Validation failed on '%s' tag", err.Tag())
			log.Error().Msg(fmt.Sprintf("Validation failed on '%s' tag", err.Tag()))
		}

		utils.NewResponse(utils.ErrValidation.Error(), http.StatusBadRequest, nil).Write(c.Writer)
		return
	}

	retried, err := h.usecase.TransactionDetailRetry(id, transactionDetail)

	if err != nil {
		utils.NewResponse(err.Error(), http.StatusInternalServerError, nil).Write(c.Writer)
		return
	}

	h.usecase.TransactionDetailSend(retried)

	utils.NewResponse("Update success", http.StatusOK, retried).Write(c.Writer)
}

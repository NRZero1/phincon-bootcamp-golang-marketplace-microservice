package impl

import (
	"donation_provider_service/internal/domain/dto/response"
	"donation_provider_service/internal/handler"
	"donation_provider_service/internal/usecase"
	"donation_provider_service/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type DonationProviderHandler struct {
	usecase usecase.DonationProviderUseCaseInterface
}

func NewBalanceHandler(usecase usecase.DonationProviderUseCaseInterface) handler.DonationProviderHandlerInterface {
	return DonationProviderHandler{
		usecase: usecase,
	}
}

func (h DonationProviderHandler) FindByID(c *gin.Context) {
	log.Trace().Msg("Entering donation provider handler find by id")

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

	foundProvider, errFound := h.usecase.FindByID(id)

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
		Data:       foundProvider,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h DonationProviderHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering donation provider get all handler")

	allProvider := h.usecase.GetAll()

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       allProvider,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Balance fetched and returning json")
	c.JSON(http.StatusOK, response)
}

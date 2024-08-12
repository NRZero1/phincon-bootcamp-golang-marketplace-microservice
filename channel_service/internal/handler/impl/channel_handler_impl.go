package impl

import (
	"channel_service/internal/domain/dto/response"
	"channel_service/internal/handler"
	"channel_service/internal/usecase"
	"channel_service/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type ChannelHandler struct {
	usecase usecase.ChannelUseCaseInterface
}

func NewChannelHandler(usecase usecase.ChannelUseCaseInterface) handler.ChannelHandlerInterface {
	return ChannelHandler{
		usecase: usecase,
	}
}

func (h ChannelHandler) FindById(c *gin.Context) {
	log.Trace().Msg("Entering courier handler find by id")

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

	foundChannel, errFound := h.usecase.FindById(id)

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
		Data:       foundChannel,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Channel fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h ChannelHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering courier get all handler")

	allChannels := h.usecase.GetAll()

	response := response.GlobalResponse{
		Message:    "OK",
		StatusCode: http.StatusOK,
		Data:       allChannels,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Channel fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h ChannelHandler) FindByName(c *gin.Context) {
	log.Trace().Msg("Entering courier handler find by name")

	name := c.Query("name")
	log.Debug().Str("Received name is: ", name)

	foundChannel, errFound := h.usecase.FindByName(name)

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
		Data:       foundChannel,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("Channel fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h ChannelHandler) GetAllOrFindByName(c *gin.Context) {
	name := c.Query("name")

	if name != "" {
		h.FindByName(c)
	} else {
		h.GetAll(c)
	}
}

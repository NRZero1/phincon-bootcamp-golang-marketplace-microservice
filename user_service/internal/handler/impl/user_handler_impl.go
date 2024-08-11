package impl

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"user_service/internal/domain/dto/request"
	"user_service/internal/domain/dto/response"
	"user_service/internal/handler"
	"user_service/internal/usecase"
	"user_service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type UserHandler struct {
	usecase usecase.UserUseCaseInterface
}

func NewUserHandler(usecase usecase.UserUseCaseInterface) (handler.UserHandlerInterface) {
	return UserHandler {
		usecase: usecase,
	}
}

func (h UserHandler) Save(c *gin.Context) {
	log.Trace().Msg("Entering user handler save")

	var register request.Register

	log.Trace().Msg("Decoding json")
	err := c.ShouldBindJSON(&register)

	if err != nil {
		log.Trace().Msg("JSON decode error")
		log.Error().Str("Error message: ", err.Error())
		response := response.GlobalResponse {
			Message: utils.ErrDecode.Error(),
			StatusCode: http.StatusBadRequest,
			Data: nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Trace().Msg("Validating user input")
	errValidate := utils.ValidateStruct(&register)
	sevenOrMore, number, upper := utils.VerifyPassword(register.Password)
	var errPassword []string

	if !sevenOrMore {
		errPassword = append(errPassword, "seven or more")
	}

	if !number {
		errPassword = append(errPassword, "number")
	}

	if !upper {
		errPassword = append(errPassword, "upper")
	}

	if errValidate != nil || len(errPassword) > 0 {
		log.Trace().Msg("Validation error")
		errors := make(map[string]string)
		log.Trace().Msg("User input error")
		if errValidate != nil {
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
		}

		if len(errPassword) > 0 {
			errors["Password"] = fmt.Sprintf("Password is not valid on [%v] validation", errPassword)
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.WriteHeader(http.StatusBadRequest)

		response := response.GlobalResponse {
			Message: utils.ErrValidation.Error(),
			StatusCode: http.StatusBadRequest,
			Data: errors,
		}

        c.JSON(http.StatusBadRequest, response)
		return
	}

	log.Debug().
		Str("Username: ", register.Username).
		Float64("Balance: ", register.Balance).
		Msg("Continuing user save process")


	savedUser, errSave := h.usecase.Save(register)

	if errSave != nil {
		log.Trace().Msg("Checking error cause")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse

		if errors.Is(errSave, utils.ErrHash) {
			log.Error().Msg(fmt.Sprintf("Error hashing password with message: %s", errSave.Error()))

			c.Writer.WriteHeader(http.StatusInternalServerError)

			response := response.GlobalResponse {
				Message: errSave.Error(),
				StatusCode: http.StatusInternalServerError,
				Data: nil,
			}
			c.JSON(http.StatusInternalServerError, response)
			return
		} else {
			log.Trace().Msg("Save error")
			log.Error().Str("Error message: ", errSave.Error())
			c.Writer.WriteHeader(http.StatusBadRequest)

			resp = response.GlobalResponse {
				Message: errSave.Error(),
				StatusCode: http.StatusBadRequest,
				Data: nil,
			}
			c.JSON(http.StatusBadRequest, resp)
		}
		return
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusCreated)

	response := response.GlobalResponse {
		Message: "Created",
		StatusCode: http.StatusCreated,
		Data: savedUser,
	}

	log.Info().Msg("User created successfully and returning json")
	c.JSON(http.StatusCreated, response)
}

func (h UserHandler) FindById(c *gin.Context) {
	log.Trace().Msg("Entering user handler find by id")

	idString := c.Param("id")
	log.Debug().Str("Received Id is: ", idString)

	log.Trace().Msg("Trying to convert id in string to int")
	id, errConv := strconv.Atoi(idString)

	if errConv != nil {
		log.Trace().Msg("Error happens when converting id string to int")
		log.Error().Str("Error message: ", errConv.Error())
		response := response.GlobalResponse {
			Message: utils.ErrPathVar.Error(),
			StatusCode: http.StatusBadRequest,
			Data: nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	foundUser, errFound := h.usecase.FindById(id)

	if errFound != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse

		log.Trace().Msg("Fetch error")
		log.Error().Str("Error message: ", errFound.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)

		resp = response.GlobalResponse {
			Message: errFound.Error(),
			StatusCode: http.StatusBadRequest,
			Data: "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	response := response.GlobalResponse {
		Message: "OK",
		StatusCode: http.StatusOK,
		Data: foundUser,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("User fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) GetAll(c *gin.Context) {
	log.Trace().Msg("Entering user get all handler")

	allUsers := h.usecase.GetAll()

	response := response.GlobalResponse {
		Message: "OK",
		StatusCode: http.StatusOK,
		Data: allUsers,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("User fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) Login(c *gin.Context) {
	var loginRequest request.Login

	err := c.ShouldBindJSON(&loginRequest)

	if err != nil {
		log.Trace().Msg("JSON decode error")
		log.Error().Str("Error message: ", err.Error())
		response := response.GlobalResponse {
			Message: err.Error(),
			StatusCode: http.StatusBadRequest,
			Data: "",
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	foundUser, err := h.usecase.FindByUsernameLogin(loginRequest.Username)

	if err != nil {
		log.Trace().Msg("Found user error")
		log.Error().Str("Error message: ", err.Error())
		response := response.GlobalResponse {
			Message: err.Error(),
			StatusCode: http.StatusNotFound,
			Data: nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusNotFound)
		c.JSON(http.StatusNotFound, response)
		return
	}

	isSame := utils.CheckPasswordHash(loginRequest.Password, foundUser.Password)

	if !isSame {
		log.Trace().Msg("Password mismatch error")
		log.Error().Str("Error message: ", errors.New("wrong password").Error())
		response := response.GlobalResponse {
			Message: utils.ErrWrongPass.Error(),
			StatusCode: http.StatusUnauthorized,
			Data: nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		c.JSON(http.StatusNotFound, response)
		return
	}

	tokenString, err := utils.GenerateJwtToken(foundUser)

	if err != nil {
		log.Trace().Msg("Error creating signature")
		log.Error().Str("Error message: ", err.Error())
		response := response.GlobalResponse {
			Message: err.Error(),
			StatusCode: http.StatusInternalServerError,
			Data: nil,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.WriteHeader(http.StatusInternalServerError)
		c.JSON(http.StatusNotFound, response)
		return
	}

	// c.SetCookie("token", tokenString, time.Now().Add(time.Minute * 1).Second(), "", "", false, true)
	response := response.GlobalResponse {
		Message: "OK",
		StatusCode: http.StatusInternalServerError,
		Data: tokenString,
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) FindByUsername(c *gin.Context) {
	log.Trace().Msg("Entering user handler find by username")

	username := c.Query("username")
	log.Debug().Str("Received username is: ", username)

	foundUser, errFound := h.usecase.FindByUsername(username)

	if errFound != nil {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")

		var resp response.GlobalResponse
		log.Trace().Msg("Fetch error")
		log.Error().Str("Error message: ", errFound.Error())
		c.Writer.WriteHeader(http.StatusBadRequest)

		resp = response.GlobalResponse {
			Message: errFound.Error(),
			StatusCode: http.StatusBadRequest,
			Data: "",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	response := response.GlobalResponse {
		Message: "OK",
		StatusCode: http.StatusOK,
		Data: foundUser,
	}

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	log.Info().Msg("User fetched and returning json")
	c.JSON(http.StatusOK, response)
}

func (h UserHandler) GetAllOrFindByName(c *gin.Context) {
	username := c.Query("username")

	if username != "" {
		h.FindByUsername(c)
	} else {
		h.GetAll(c)
	}
}

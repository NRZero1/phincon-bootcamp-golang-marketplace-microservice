package impl

import (
	"user_service/internal/domain/dto/request"
	"user_service/internal/domain/dto/response"
	"user_service/internal/repository"
	"user_service/internal/usecase"
	"user_service/internal/utils"

	"github.com/rs/zerolog/log"
)

type UserUseCase struct {
	repo repository.UserRepositoryInterface
}

func NewUserUseCase(repo repository.UserRepositoryInterface) (usecase.UserUseCaseInterface) {
	return UserUseCase{
		repo: repo,
	}
}

func (uc UserUseCase) Save(user request.Register) (response.UserResponse, error) {
	log.Trace().Msg("Entering user usecase save")
	hashedPass, err := utils.HashPassword(user.Password)
	if err != nil {
		return response.UserResponse{}, utils.ErrHash
	}

	user.Password = hashedPass

	resp, errSave := uc.repo.Save(&user)
	return resp, errSave
}

func (uc UserUseCase) FindById(id int) (response.UserResponse, error) {
	log.Trace().Msg("Entering user usecase find by id")
	return uc.repo.FindByID(id)
}

func (uc UserUseCase) GetAll() ([]response.UserResponse) {
	log.Trace().Msg("Entering user usecase get all")
	return uc.repo.GetAll()
}

func (uc UserUseCase) ReduceBalance(id int, amount float64) (response.UserResponse, error) {
	log.Trace().Msg("Entering user usecase reduce balance")
	return uc.repo.ReduceBalance(id, amount)
}

func (uc UserUseCase) FindByUsernameLogin(username string) (response.LoginResponse, error) {
	log.Trace().Msg("Enter user usecase FindByUsernameLogin")
	return uc.repo.FindByUsernameLogin(username)
}

func (uc UserUseCase) FindByUsername(username string) (response.UserResponse, error) {
	log.Trace().Msg("Enter user usecase FindByUsername")
	return uc.repo.FindByUsername(username)
}

func (uc UserUseCase) SetPackage(userID int, packageID int) (response.UserResponse, error) {
	log.Trace().Msg("Enter user use case SetPackage")
	return uc.repo.SetPackage(userID, packageID)
}

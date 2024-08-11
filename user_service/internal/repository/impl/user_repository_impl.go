package impl

import (
	"fmt"
	"sync"
	"user_service/internal/domain"
	"user_service/internal/domain/dto/request"
	"user_service/internal/domain/dto/response"
	"user_service/internal/repository"
	"user_service/internal/utils"

	"github.com/rs/zerolog/log"
)

type UserRepository struct {
	mtx sync.Mutex
	users map[int]domain.User
	nextId int
}

func NewUserRepository() repository.UserRepositoryInterface {
	repo := &UserRepository {
		users: map[int]domain.User{},
		nextId: 1,
	}
	repo.initData()
	return repo
}

func (repo *UserRepository) initData() {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.users[1] = domain.User{
		UserID: 1,
		Username: "Daniel",
		Password: "test123",
		Balance: 1000000,
		PackageID: 0,
	}

	repo.users[2] = domain.User{
		UserID: 2,
		Username: "Ahmad",
		Password: "test123",
		Balance: 1000000,
		PackageID: 0,
	}
}

func (repo *UserRepository) Save(registerUser *request.Register) (response.UserResponse, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository save")
	log.Trace().Msg("Attempting to save new user")

	user := domain.User {
		UserID: repo.nextId,
		Username: registerUser.Username,
		Password: registerUser.Password,
		Balance: registerUser.Balance,
	}
	repo.users[user.UserID] = user
	log.Trace().Msg("New user saved")

	temp := repo.users[user.UserID]
	resp := response.UserResponse {
		UserID: temp.UserID,
		Username: temp.Username,
		Balance: temp.Balance,
	}
	return resp, nil
}

func (repo *UserRepository) FindByID(id int) (response.UserResponse, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository find by id")
	log.Trace().Msg("Attempting to fetch user")
	if foundUser, exists := repo.users[id]; exists {
		log.Trace().Msg("Fetching completed")
		temp := response.UserResponse {
			UserID: foundUser.UserID,
			Username: foundUser.Username,
			Balance: foundUser.Balance,
			PackageID: foundUser.PackageID,
		}
		return temp, nil
	}
	log.Error().Msg(fmt.Sprintf("User with ID %d not found", id))
	return response.UserResponse{}, utils.NewErrFindUserById(id)
}

func (repo *UserRepository) GetAll() ([]response.UserResponse) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository get all")
	log.Trace().Msg("Attempting to fetch user")
	listOfUsers := make([]response.UserResponse, 0, len(repo.users))

	for _, v := range repo.users {
		temp := response.UserResponse {
			UserID: v.UserID,
			Username: v.Username,
			Balance: v.Balance,
			PackageID: v.PackageID,
		}
		listOfUsers = append(listOfUsers, temp)
	}

	log.Trace().Msg("Fetching completed")
	return listOfUsers
}

func (repo *UserRepository) ReduceBalance(id int, amount float64) (response.UserResponse, error) {
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository reduce balance")
	log.Trace().Msg("Attempting to reduce user balance")
	log.Trace().Msg("Attempting to fetch user")
	if foundUser, exists := repo.users[id]; exists {
		log.Trace().Msg("Fetching completed")
		repo.mtx.Lock()
		log.Debug().Msg(fmt.Sprintf("Balance before reduced: %f", foundUser.Balance))
		foundUser.Balance = foundUser.Balance - amount

		repo.users[foundUser.UserID] = foundUser
		log.Debug().Msg(fmt.Sprintf("Balance after reduced: %f", foundUser.Balance))
		log.Trace().Msg("Balance reduced successfully")
		temp := response.UserResponse {
			UserID: foundUser.UserID,
			Username: foundUser.Username,
			Balance: foundUser.Balance,
		}
		return temp, nil
	} else {
		log.Error().Msg(fmt.Sprintf("User with ID %d not found", id))
		return response.UserResponse{}, utils.NewErrFindUserById(id)
	}
}

func (repo *UserRepository) FindByUsernameLogin(username string) (response.LoginResponse, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository find by username")
	for _, v := range repo.users {
		if v.Username == username {
			temp := response.LoginResponse {
				UserID: v.UserID,
				Username: v.Username,
				Password: v.Password,
			}
			return temp, nil
		}
	}
	return response.LoginResponse{}, utils.NewErrFindByUsername(username)
}

func (repo *UserRepository) FindByUsername(username string) (response.UserResponse, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository find by username")
	for _, v := range repo.users {
		if v.Username == username {
			temp := response.UserResponse {
				UserID: v.UserID,
				Username: v.Username,
				Balance: v.Balance,
				PackageID: v.PackageID,
			}
			return temp, nil
		}
	}
	return response.UserResponse{}, utils.NewErrFindByUsername(username)
}

func (repo *UserRepository) SetPackage(userID int, packageID int) (response.UserResponse, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside user repository set package")
	if foundUser, exists := repo.users[userID]; exists {
		log.Trace().Msg("Fetching completed")
		foundUser.PackageID = packageID
		repo.users[userID] = foundUser

		temp := response.UserResponse {
			UserID: foundUser.UserID,
			Username: foundUser.Username,
			Balance: foundUser.Balance,
			PackageID: foundUser.PackageID,
		}
		return temp, nil
	}
	return response.UserResponse{}, utils.NewErrFindUserById(userID)
}

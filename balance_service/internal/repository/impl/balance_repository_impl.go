package impl

import (
	"balance_service/internal/domain"
	"balance_service/internal/repository"
	"balance_service/internal/utils"
	"fmt"

	"github.com/rs/zerolog/log"
)

type BalanceRepository struct {
	balances map[int]domain.Balance
	nextId int
}

func NewBalanceRepository() repository.BalanceRepositoryInterface {
	repo := &BalanceRepository{
		balances: map[int]domain.Balance{},
		nextId: 1,
	}
	repo.initData()
	return repo
}

func (repo *BalanceRepository) initData() {
	repo.balances[1] = domain.Balance{
		UserID: 1,
		Balance: 10000000,
	}

	repo.balances[2] = domain.Balance{
		UserID: 2,
		Balance: 2000000000,
	}
}

func (repo *BalanceRepository) FindByID(id int) (domain.Balance, error) {
	log.Trace().Msg("Inside balance repository find by id")
	log.Trace().Msg("Attempting to fetch balance")
	if foundBalance, exists := repo.balances[id]; exists {
		log.Trace().Msg("Fetching completed")
		return foundBalance, nil
	}
	log.Error().Msg(fmt.Sprintf("Balance with ID %d not found", id))
	return domain.Balance{}, utils.NewErrFindById(id)
}

func (repo *BalanceRepository) GetAll() []domain.Balance {
	log.Trace().Msg("Inside balance repository get all")
	log.Trace().Msg("Attempting to fetch balances")
	listOfBalances := make([]domain.Balance, 0, len(repo.balances))

	for _, v := range repo.balances {
		temp := domain.Balance{
			UserID: v.UserID,
			Balance: v.Balance,
		}
		listOfBalances = append(listOfBalances, temp)
	}

	log.Trace().Msg("Fetching completed")
	return listOfBalances
}

func (repo *BalanceRepository) Deduct(id int, amount float64) (error) {

	log.Trace().Msg("Inside balance repository deduct")

	log.Trace().Msg("Attempting to reduce balance")
	foundBalance, err := repo.FindByID(id)

	if err != nil {
		return err
	}

	log.Debug().Msg(fmt.Sprintf("Balance before reduced: %f", foundBalance.Balance))
	foundBalance.Balance = foundBalance.Balance - amount
	repo.balances[foundBalance.UserID] = foundBalance
	log.Debug().Msg(fmt.Sprintf("Balance after reduced: %f", foundBalance.Balance))
	log.Info().Msg("Balance reduced successfully")
	return nil
}

func (repo *BalanceRepository) AddBalance(id int, amount float64) (error) {

	log.Trace().Msg("Inside balance repository add balance")

	log.Trace().Msg("Attempting to add balance")
	foundBalance, err := repo.FindByID(id)

	if err != nil {
		return err
	}

	log.Debug().Msg(fmt.Sprintf("Balance before added: %f", foundBalance.Balance))
	foundBalance.Balance = foundBalance.Balance + amount
	repo.balances[foundBalance.UserID] = foundBalance
	log.Debug().Msg(fmt.Sprintf("Balance after added: %f", foundBalance.Balance))
	log.Info().Msg("Balance added successfully")
	return nil
}

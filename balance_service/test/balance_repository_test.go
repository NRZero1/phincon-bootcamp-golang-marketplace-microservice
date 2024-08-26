package impl_test

import (
	"balance_service/internal/repository/impl"
	"balance_service/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBalanceRepository(t *testing.T) {
	repo := impl.NewBalanceRepository()

	assert.NotNil(t, repo)
	assert.Equal(t, 2, len(repo.GetAll()), "The repository should have initialized with 2 balances")
}

func TestFindByID(t *testing.T) {
	repo := impl.NewBalanceRepository()

	t.Run("Existing ID", func(t *testing.T) {
		balance, err := repo.FindByID(1)
		assert.NoError(t, err)
		assert.Equal(t, 10000000.0, balance.Balance, "The balance for user 1 should be 10000000")
	})

	t.Run("Non-existing ID", func(t *testing.T) {
		_, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Equal(t, utils.NewErrFindById(999), err, "The error should be for non-existing ID")
	})
}

func TestGetAll(t *testing.T) {
	repo := impl.NewBalanceRepository()

	balances := repo.GetAll()
	assert.Equal(t, 2, len(balances), "There should be 2 balances in the repository")
	assert.Equal(t, 10000000.0, balances[0].Balance, "The balance for user 1 should be 10000000")
	assert.Equal(t, 2000000000.0, balances[1].Balance, "The balance for user 2 should be 2000000000")
}

func TestDeduct(t *testing.T) {
	repo := impl.NewBalanceRepository()

	t.Run("Successful Deduction", func(t *testing.T) {
		err := repo.Deduct(1, 5000000)
		assert.NoError(t, err)

		balance, _ := repo.FindByID(1)
		assert.Equal(t, 5000000.0, balance.Balance, "The balance for user 1 should be reduced to 5000000")
	})

	t.Run("Deduction with Non-existing ID", func(t *testing.T) {
		err := repo.Deduct(999, 5000000)
		assert.Error(t, err)
		assert.Equal(t, utils.NewErrFindById(999), err, "The error should be for non-existing ID")
	})
}

func TestAddBalance(t *testing.T) {
	repo := impl.NewBalanceRepository()

	t.Run("Successful Addition", func(t *testing.T) {
		err := repo.AddBalance(1, 5000000)
		assert.NoError(t, err)

		balance, _ := repo.FindByID(1)
		assert.Equal(t, 15000000.0, balance.Balance, "The balance for user 1 should be increased to 15000000")
	})

	t.Run("Addition with Non-existing ID", func(t *testing.T) {
		err := repo.AddBalance(999, 5000000)
		assert.Error(t, err)
		assert.Equal(t, utils.NewErrFindById(999), err, "The error should be for non-existing ID")
	})
}

package impl

import (
	"donation_provider_service/internal/domain"
	"donation_provider_service/internal/repository"
	"donation_provider_service/internal/utils"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type DonationProviderRepository struct {
	mtx sync.Mutex
	donations map[int]domain.Provider
	nextId int
}

func NewDonationProviderRepository() repository.DonationProviderRepositoryInterface {
	repo := &DonationProviderRepository{
		donations: map[int]domain.Provider{},
		nextId: 1,
	}
	repo.initData()
	return repo
}

func (repo *DonationProviderRepository) initData() {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.donations[1] = domain.Provider{
		ProviderID: 1,
		Name: "Test provider 1",
	}

	repo.donations[2] = domain.Provider{
		ProviderID: 2,
		Name: "Test provider 2",
	}
}

func (repo *DonationProviderRepository) FindByID(id int) (domain.Provider, error) {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside donation repository find by id")
	log.Trace().Msg("Attempting to fetch donation")
	if foundDonation, exists := repo.donations[id]; exists {
		log.Trace().Msg("Fetching completed")
		return foundDonation, nil
	}
	log.Error().Msg(fmt.Sprintf("Donation with ID %d not found", id))
	return domain.Provider{}, utils.NewErrFindById(id)
}

func (repo *DonationProviderRepository) GetAll() []domain.Provider {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	log.Trace().Msg("Inside donation repository get all")
	log.Trace().Msg("Attempting to fetch donations")
	listOfDonations := make([]domain.Provider, 0, len(repo.donations))

	for _, v := range repo.donations {
		temp := domain.Provider{
			ProviderID: v.ProviderID,
			Name: v.Name,
		}
		listOfDonations = append(listOfDonations, temp)
	}

	log.Trace().Msg("Fetching completed")
	return listOfDonations
}

package impl

import (
	"channel_service/internal/domain"
	"channel_service/internal/repository"
	"channel_service/internal/utils"
	"fmt"
	"sync"

	"github.com/rs/zerolog/log"
)

type ChannelRepository struct {
	mtx sync.Mutex
	channels map[int]domain.Channel
	nextId int
}

func NewChannelRepository() repository.ChannelRepositoryInterface {
	repo := &ChannelRepository{
		channels: map[int]domain.Channel{},
		nextId: 1,
	}
	repo.initData()
	return repo
}

func (repo *ChannelRepository) initData() {
	repo.mtx.Lock()
	defer repo.mtx.Unlock()

	repo.channels[1] = domain.Channel{
		ChannelID:        1,
		UserID:           1,
		ChannelName:      "Test Channel 1",
		MembershipPrice:  10000,
	}

	repo.channels[2] = domain.Channel{
		ChannelID:        2,
		UserID:           2,
		ChannelName:      "Test Channel 2",
		MembershipPrice:  10000,
	}
}

func (repo *ChannelRepository) Save(channel *domain.Channel) (domain.Channel, error) {

	log.Trace().Msg("Inside channel repository save")
	log.Trace().Msg("Attempting to save new channel")

	channel.ChannelID = repo.nextId
	repo.channels[channel.ChannelID] = *channel
	log.Trace().Msg("New channel saved")

	return repo.channels[channel.ChannelID], nil
}

func (repo *ChannelRepository) FindById(id int) (domain.Channel, error) {

	log.Trace().Msg("Inside channel repository find by id")
	log.Trace().Msg("Attempting to fetch channel")
	if foundChannel, exists := repo.channels[id]; exists {
		log.Trace().Msg("Fetching completed")
		log.Debug().Msgf("Found Channel Value: %+v", foundChannel)
		return foundChannel, nil
	}
	log.Error().Msg(fmt.Sprintf("Channel with ID %d not found", id))
	return domain.Channel{}, utils.NewErrFindById(id)
}

func (repo *ChannelRepository) GetAll() []domain.Channel {

	log.Trace().Msg("Inside channel repository get all")
	log.Trace().Msg("Attempting to fetch channels")
	listOfChannels := make([]domain.Channel, 0, len(repo.channels))

	for _, v := range repo.channels {
		temp := domain.Channel{
			ChannelID:       v.ChannelID,
			ChannelName:     v.ChannelName,
			MembershipPrice: v.MembershipPrice,
			UserID: v.UserID,
			Membership: v.Membership,
		}
		listOfChannels = append(listOfChannels, temp)
	}

	log.Trace().Msg("Fetching completed")
	return listOfChannels
}

func (repo *ChannelRepository) FindByName(name string) (domain.Channel, error) {

	log.Trace().Msg("Inside channel repository find by name")
	for _, v := range repo.channels {
		if v.ChannelName == name {
			temp := domain.Channel{
				ChannelID:       v.ChannelID,
				ChannelName:     v.ChannelName,
				MembershipPrice: v.MembershipPrice,
				UserID: v.UserID,
				Membership: v.Membership,
			}
			return temp, nil
		}
	}
	return domain.Channel{}, utils.NewErrFindByName(name)
}

func (repo *ChannelRepository) AddMembership(channelID int, userID int) (domain.Channel, error) {

	if foundChannel, exists := repo.channels[channelID]; exists {
		log.Trace().Msg("Fetching completed")

		foundChannel.Membership = append(foundChannel.Membership, userID)
		repo.channels[channelID] = foundChannel
		return repo.channels[channelID], nil
	}
	return domain.Channel{}, utils.NewErrFindById(channelID)
}

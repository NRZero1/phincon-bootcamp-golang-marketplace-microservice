package impl_test

import (
	"channel_service/internal/domain"
	"channel_service/internal/repository/impl"
	"channel_service/internal/utils"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannelRepository_Save(t *testing.T) {
	repo := impl.NewChannelRepository()

	channel := &domain.Channel{
		ChannelName:     "New Channel",
		MembershipPrice: 15000,
		UserID:          3,
	}

	savedChannel, err := repo.Save(channel)

	assert.NoError(t, err)
	assert.Equal(t, 3, savedChannel.ChannelID) // Changed to match `repo.nextId` increment
	assert.Equal(t, "New Channel", savedChannel.ChannelName)
	assert.Equal(t, 15000.0, savedChannel.MembershipPrice)
	assert.Equal(t, 3, savedChannel.UserID)
}

func TestChannelRepository_FindById(t *testing.T) {
	repo := impl.NewChannelRepository()

	channel, err := repo.FindById(1)

	assert.NoError(t, err)
	assert.Equal(t, 1, channel.ChannelID)
	assert.Equal(t, "Test Channel 1", channel.ChannelName)

	_, err = repo.FindById(99)
	assert.Error(t, err)
	assert.Equal(t, utils.NewErrFindById(99).Error(), err.Error())
}

func TestChannelRepository_GetAll(t *testing.T) {
	repo := impl.NewChannelRepository()

	channels := repo.GetAll()

	assert.Len(t, channels, 2)
	assert.Equal(t, 1, channels[0].ChannelID)
	assert.Equal(t, 2, channels[1].ChannelID)
}

func TestChannelRepository_FindByName(t *testing.T) {
	repo := impl.NewChannelRepository()

	channel, err := repo.FindByName("Test Channel 1")

	assert.NoError(t, err)
	assert.Equal(t, 1, channel.ChannelID)
	assert.Equal(t, "Test Channel 1", channel.ChannelName)

	_, err = repo.FindByName("Nonexistent Channel")
	assert.Error(t, err)
	assert.Equal(t, utils.NewErrFindByName("Nonexistent Channel").Error(), err.Error())
}

func TestChannelRepository_AddMembership(t *testing.T) {
	repo := impl.NewChannelRepository()

	channel, err := repo.AddMembership(1, 3)

	assert.NoError(t, err)
	assert.Equal(t, 1, channel.ChannelID)
	assert.Contains(t, channel.Membership, 3)

	_, err = repo.AddMembership(99, 3)
	assert.Error(t, err)
	assert.Equal(t, utils.NewErrFindById(99).Error(), err.Error())
}

func TestChannelRepository_Concurrency(t *testing.T) {
	repo := impl.NewChannelRepository()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			repo.Save(&domain.Channel{
				ChannelName:     "Channel" + strconv.Itoa(n), // Convert integer to string
				MembershipPrice: float64(n * 1000),
				UserID:          n,
			})
		}(i)
	}

	wg.Wait()
	channels := repo.GetAll()

	assert.Len(t, channels, 102) // Including initial 2 channels
}

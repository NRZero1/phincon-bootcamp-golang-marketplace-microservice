package repository

import (
	"channel_service/internal/repository"
	repoImplement "channel_service/internal/repository/impl"
)

var (
	ChannelRepository repository.ChannelRepositoryInterface
)

func InitRepository() {
	ChannelRepository = repoImplement.NewChannelRepository()
}

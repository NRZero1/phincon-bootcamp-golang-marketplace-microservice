package repository

import (
	"channel_service/internal/domain"
)

type ChannelRepositoryInterface interface {
	ChannelSave
	ChannelFindByName
	ChannelFindById
	ChannelGetAll
	ChannelAddMembership
}

type ChannelSave interface {
	Save(channel *domain.Channel) (domain.Channel, error)
}

type ChannelFindByName interface {
	FindByName(name string) (domain.Channel, error)
}

type ChannelFindById interface {
	FindById(id int) (domain.Channel, error)
}

type ChannelGetAll interface {
	GetAll() ([]domain.Channel)
}

type ChannelAddMembership interface {
	AddMembership(channelID int, userID int) (domain.Channel, error)
}

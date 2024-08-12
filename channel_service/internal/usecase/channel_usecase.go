package usecase

import "channel_service/internal/domain"

type ChannelUseCaseInterface interface {
	ChannelFindById
	ChannelGetAll
	ChannelFindByName
	ChannelAddMembership
}

type ChannelFindById interface {
	FindById(id int) (domain.Channel, error)
}

type ChannelGetAll interface {
	GetAll() ([]domain.Channel)
}

type ChannelFindByName interface {
	FindByName(name string) (domain.Channel, error)
}

type ChannelAddMembership interface {
	AddMembership(channelID int, userID int) (domain.Channel, error)
}

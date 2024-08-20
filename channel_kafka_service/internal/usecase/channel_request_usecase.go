package usecase

import response "channel_kafka_service/internal/domain/channel_service_response"

type ChannelRequestUseCaseInterface interface {
	ChannelGetByID
}

type ChannelGetByID interface {
	GetChannelByID(id int) (bool, int, response.Channel, error)
}

type ChannelAddMembership interface {
	AddMembership(channelID int, userID int) (bool, int, error)
}

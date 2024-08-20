package usecase

type ConsumerUseCaseInterface interface {
	ConsumerRouting
}

type ConsumerRouting interface {
	RouteMessage(message []byte)
}

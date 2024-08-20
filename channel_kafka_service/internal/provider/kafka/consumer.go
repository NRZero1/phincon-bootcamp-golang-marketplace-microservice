package kafka

import (
	"channel_kafka_service/internal/usecase"
	"channel_kafka_service/internal/utils"
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	*kafka.Reader
	consumerUseCase usecase.ConsumerUseCaseInterface
}

func NewKafkaConsumer(broker []string, groupID string, topic string, consumerUseCase usecase.ConsumerUseCaseInterface) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: broker,
			GroupID: groupID,
			Topic: topic,
			MaxBytes: 10e3, // 10KB
			StartOffset: kafka.LastOffset,
		}),
		consumerUseCase: consumerUseCase,
	}
}

func (c *KafkaConsumer) ConsumeMessage(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka consumer shutting down...")
			return nil
		default:
			message, err := c.ReadMessage(ctx)

			if err != nil {
				log.Error().Msg(fmt.Sprintf("Error when trying to read message: %s", err.Error()))
				return utils.ErrKafkaConsume
			}

			log.Info().Msg(fmt.Sprintf("Message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value)))
			c.consumerUseCase.RouteMessage(message.Value)
		}
	}
}

func (c *KafkaConsumer) Close() {
	if err := c.Reader.Close(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error hwne trying to close kafka consumer with message: %s", err.Error()))
	}
}

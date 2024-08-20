package kafka

import (
	"context"
	"fmt"
	"order_service/internal/utils"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	*kafka.Reader
}

func NewKafkaConsumer(broker []string, groupID string, topic string) *KafkaConsumer {
	return &KafkaConsumer{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: broker,
			GroupID: groupID,
			Topic: topic,
			MaxBytes: 10e3, // 10KB
			StartOffset: kafka.LastOffset,
		}),
	}
}

func (c *KafkaConsumer) ConsumeMessage() error {
	for {
		message, err := c.ReadMessage(context.Background())

		if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to read message with error message: %s", err.Error()))
			return utils.ErrKafkaConsume
		}

		log.Info().Msg(fmt.Sprintf("message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value)))
	}
}

func (c *KafkaConsumer) Close() {
	if err := c.Reader.Close(); err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to close kafka consumer with message: %s", err.Error()))
	}
}

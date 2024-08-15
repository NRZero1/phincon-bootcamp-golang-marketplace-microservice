package kafka

import (
	"context"
	"fmt"
	"orchestration_service/internal/utils"

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
		}),
	}
}

func (c *KafkaConsumer) ConsumeMessage() error {
	for {
		message, err := c.ReadMessage(context.Background())

		if err != nil {
			log.Error().Msg(fmt.Sprintf("%s", utils.NewErrKafkaConsume(err.Error())))
			return err
		}

		log.Info().Msg(fmt.Sprintf("message at offset %d: %s = %s\n", message.Offset, string(message.Key), string(message.Value)))
	}
}

func (c *KafkaConsumer) Close() {
	if err := c.Reader.Close(); err != nil {
		log.Error().Msg(fmt.Sprintf("%s", utils.NewErrKafkaReaderClose(err.Error())))
	}
}
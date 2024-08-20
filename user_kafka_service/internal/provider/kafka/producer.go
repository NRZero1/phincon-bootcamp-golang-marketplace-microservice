package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"user_kafka_service/internal/utils"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
    *kafka.Writer
    *kafka.Conn
}

func NewKafkaProducer(brokers []string, topic string, partitions int, replicationFactor int) (*KafkaProducer, error) {
    conn, err := kafka.Dial("tcp", brokers[0])
    if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to create connection to kafka with message: %s", err.Error()))
        return nil, utils.ErrKafkaProducer
    }

    topicExists, err := checkTopicExists(conn, topic)
    if err != nil {
		log.Error().Msg(fmt.Sprintf("Error when trying to check if topic is exist with message: %s", err.Error()))
        return nil, utils.ErrKafkaProducer
    }

    if !topicExists {
        err = createTopic(conn, topic, partitions, replicationFactor)
        if err != nil {
			log.Error().Msg(fmt.Sprintf("Error when trying to create topic with message: %s", err.Error()))
            return nil, utils.ErrKafkaProducer
        }
    }

    return &KafkaProducer{
        Writer: kafka.NewWriter(kafka.WriterConfig{
            Brokers:  brokers,
            Topic:    topic,
            Balancer: &kafka.LeastBytes{},
        }),
        Conn: conn,
    }, nil
}

func checkTopicExists(conn *kafka.Conn, topic string) (bool, error) {
    partitions, err := conn.ReadPartitions(topic)
    if err != nil {
        return false, err
    }

    return len(partitions) > 0, nil
}

func createTopic(conn *kafka.Conn, topic string, partitions int, replicationFactor int) error {
    topicConfigs := []kafka.TopicConfig{
        {
            Topic:             topic,
            NumPartitions:     partitions,
            ReplicationFactor: replicationFactor,
        },
    }

    err := conn.CreateTopics(topicConfigs...)
    if err != nil {
        log.Printf("Failed to create topic: %v", err)
        return err
    }

    log.Printf("Topic '%s' created successfully", topic)
    return nil
}

func (kp *KafkaProducer) ProduceMessage(key string, value any) error {
	valueByte, err := json.MarshalIndent(value, "", "	")

	if err != nil {
		log.Error().Msg(fmt.Sprintf("Failed to marshall value with message: %s", err.Error()))
		return utils.ErrMarshal
	}

    msg := kafka.Message{
        Key:   []byte(key),
        Value: valueByte,
    }

    ctx := context.Background()
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    errWrite := kp.Writer.WriteMessages(ctx, msg)

    if errWrite != nil {
        log.Printf("Failed to write message to Kafka: %v", err)
        return err
    }

    log.Info().Msg(fmt.Sprintf("Produced message to Kafka: key=%s, value=%+v", key, value))
    return nil
}

func (kp *KafkaProducer) Close() {
    if err := kp.Writer.Close(); err != nil {
        log.Printf("Failed to close Kafka writer: %v", err)
    }
    if err := kp.Conn.Close(); err != nil {
        log.Printf("Failed to close Kafka admin connection: %v", err)
    }
}

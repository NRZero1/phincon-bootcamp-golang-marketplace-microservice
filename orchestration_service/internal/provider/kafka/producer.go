package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
    *kafka.Writer
    *kafka.Conn
}

func NewKafkaProducer(brokers []string, topic string, partitions int, replicationFactor int) (*KafkaProducer, error) {
    conn, err := kafka.Dial("tcp", brokers[0])
    if err != nil {
        return nil, err
    }

    topicExists, err := checkTopicExists(conn, topic)
    if err != nil {
        return nil, err
    }

    if !topicExists {
        err = createTopic(conn, topic, partitions, replicationFactor)
        if err != nil {
            return nil, err
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

func (kp *KafkaProducer) ProduceMessage(key, value string) error {
    msg := kafka.Message{
        Key:   []byte(key),
        Value: []byte(value),
    }

    ctx := context.Background()
    ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
    defer cancel()

    err := kp.Writer.WriteMessages(ctx, msg)
    if err != nil {
        log.Printf("Failed to write message to Kafka: %v", err)
        return err
    }

    log.Printf("Produced message to Kafka: key=%s, value=%s", key, value)
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

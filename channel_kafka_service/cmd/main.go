package main

import (
	"channel_kafka_service/internal/provider/kafka"
	"channel_kafka_service/internal/provider/usecase"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

var brokers []string
var groupID string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	brokers = append(brokers, os.Getenv("KAFKA_BROKER"))
	groupID = os.Getenv("CONSUME_GROUPID")
	usecase.InitUseCase()
}

func main() {
	topicsStr := os.Getenv("KAFKA_TOPICS")
	if topicsStr == "" {
		log.Fatal().Msg("KAFKA_TOPICS environment variable is not set")
	}

	// Split the topics string into a slice
	topics := strings.Split(topicsStr, ",")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(len(topics))

	log.Info().Msg("Orchestration service is running and consuming messages...")

	// Create consumers for each topic and run them in separate goroutines
	for _, topic := range topics {
		consumer := kafka.NewKafkaConsumer(brokers, groupID, topic, usecase.ConsumerUseCase)

		go func(c *kafka.KafkaConsumer, t string) {
			defer wg.Done()
			log.Info().Msg(fmt.Sprintf("Starting consumer for topic: %s", t))
			if err := c.ConsumeMessage(ctx); err != nil {
				log.Error().Msg(fmt.Sprintf("Failed to consume Kafka messages from topic %s: %v", t, err))
			}
		}(consumer, topic)
	}

	// Wait for a termination signal
	<-sigchan
	log.Info().Msg("Shutting down Orchestration service...")

	// Signal all consumers to stop
	cancel()

	// Wait for all consumers to finish
	wg.Wait()

	log.Info().Msg("Service shut down successfully.")
}

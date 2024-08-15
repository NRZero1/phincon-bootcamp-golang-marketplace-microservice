package main

import (
	"database/sql"
	"fmt"
	"orchestration_service/internal/provider/db"
	"orchestration_service/internal/provider/kafka"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

var brokers []string
var groupID string
var database *sql.DB

func initDB() *sql.DB {
	db, err := db.NewConnection(os.Getenv("DB")).GetConnection(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	if err != nil {
		return nil
	}

	return db
}

func init() {
	brokers = append(brokers, os.Getenv("KAFKA_BROKER"))
	groupID = os.Getenv("CONSUME_GROUPID")
	database = initDB()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}

	topicsStr := os.Getenv("KAFKA_TOPICS")
	if topicsStr == "" {
		log.Fatal().Msg("KAFKA_TOPICS environment variable is not set")
	}

	// Split the topics string into a slice
	topics := strings.Split(topicsStr, ",")

	sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

    log.Info().Msg("Orchestration service is running and consuming messages...")

    for {
        select {
        case <-sigchan:
            log.Info().Msg("Shutting down Orchestration service...")
			errDbClose := database.Close()
			if errDbClose != nil {
				log.Fatal().Msg(fmt.Sprintf("Database shutdown error: %s", errDbClose.Error()))
			}
            return
        default:
            for _, topic := range topics {
				consumer := kafka.NewKafkaConsumer(brokers, groupID, topic)
				if err != nil {
					log.Fatal().Msg(fmt.Sprintf("Failed to create Kafka consumer: %v", err))
				}

				go func(c *kafka.KafkaConsumer, t string) {
					log.Printf("Starting consumer for topic: %s", t)
					if err := c.ConsumeMessage(); err != nil {
						log.Fatal().Msg(fmt.Sprintf("Failed to consume Kafka messages from topic %s: %v", t, err))
					}
				}(consumer, topic)
			}
        }
    }

}
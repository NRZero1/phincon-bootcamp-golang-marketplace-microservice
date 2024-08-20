package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"orchestration_service/internal/provider/db"
	"orchestration_service/internal/provider/handler"
	"orchestration_service/internal/provider/kafka"
	"orchestration_service/internal/provider/repository"
	"orchestration_service/internal/provider/routes"
	"orchestration_service/internal/provider/usecase"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	brokers = append(brokers, os.Getenv("KAFKA_BROKER"))
	groupID = os.Getenv("CONSUME_GROUPID")
	database = initDB()
	repository.InitRepository(database)
	usecase.InitUseCase()
	handler.InitHandler()
}

func main() {
	router := gin.Default()

	globalGroup := router.Group("")
	{
		routes.TransactionRoutes(globalGroup.Group("/transaction"), handler.TransactionHandler)
	}

	server := &http.Server{
		Addr:    "localhost:8091",
		Handler: router,
	}

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
	wg.Add(len(topics) + 1)

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

	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msg(fmt.Sprintf("Failed to start Gin server: %s", err.Error()))
		}
	}()

	// Wait for a termination signal
	<-sigchan
	log.Info().Msg("Shutting down Orchestration service...")

	// Signal all consumers to stop
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Msg(fmt.Sprintf("Gin server shutdown error: %s", err.Error()))
	}

	// Wait for all consumers to finish
	wg.Wait()

	// Close the database connection
	errDbClose := database.Close()
	if errDbClose != nil {
		log.Fatal().Msg(fmt.Sprintf("Database shutdown error: %s", errDbClose.Error()))
	}
	log.Info().Msg("Service shut down successfully.")
}

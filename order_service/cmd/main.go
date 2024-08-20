package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"order_service/internal/provider/db"
	"order_service/internal/provider/handler"
	"order_service/internal/provider/repository"
	"order_service/internal/provider/routes"
	"order_service/internal/provider/usecase"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

var database *sql.DB

func initDB() *sql.DB {
	log.Debug().Msgf("DB: %s", os.Getenv("DB"))
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
	database = initDB()
	repository.InitRepository(database)
	usecase.InitUseCase()
	handler.InitHandler()
}

func main() {
	router := gin.Default()

	globalRoutesGroup := router.Group("")
	{
		routes.UserRoutes(globalRoutesGroup.Group("/user"), "http://localhost:8080")
		routes.OrderRoutes(globalRoutesGroup.Group("/order"), handler.OrderHandler)
	}

	router.POST("/order", )

	server := &http.Server{
		Addr:    "localhost:8090",
		Handler: router,
	}

	go func() {
		fmt.Println("Server is running in port ", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msg(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    log.Info().Msg("Shutting down the server...")

    // Set a timeout for shutdown (for example, 5 seconds).
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatal().Msg(fmt.Sprintf("Server shutdown error: %v", err))
    }
    log.Info().Msg("Server gracefully stopped")
}

package main

import (
	"context"
	"database/sql"
	"fmt"
	"gateway/internal/provider/db"
	"gateway/internal/provider/handler"
	"gateway/internal/provider/repository"
	"gateway/internal/provider/routes"
	"gateway/internal/provider/usecase"
	"net/http"
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
	db, err := db.NewConnection(os.Getenv("DB")).GetConnection(os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	if err != nil {
		return nil
	}

	return db
}

func init() {
	database = initDB()
	repository.InitRepository(database)
	usecase.InitUseCase()
	handler.InitHandler()
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("Error loading .env file")
	}
	router := gin.Default()

	globalRoutesGroup := router.Group("")
	{
		routes.UserRoutes(globalRoutesGroup.Group("/user"), "http://localhost:8080")
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

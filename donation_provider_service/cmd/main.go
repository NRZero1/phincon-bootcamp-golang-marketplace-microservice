package main

import (
	"context"
	"donation_provider_service/internal/provider/handler"
	"donation_provider_service/internal/provider/repository"
	"donation_provider_service/internal/provider/routes"
	"donation_provider_service/internal/provider/usecase"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	// zerolog.SetGlobalLevel(zerolog.InfoLevel)
	router := gin.New()

	// middlewares := middleware.CreateStack(
	// 	middleware.Logging,
	// )

	repository.InitRepository()
	usecase.InitUseCase()
	handler.InitHandler()

	router.Use(gin.Recovery())

	globalGroup := router.Group("")
	{
		routes.DonationRoutes(globalGroup.Group("/donation"), handler.DonationProviderHandler)
	}

	server := &http.Server{
		Addr:    "localhost:8083",
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

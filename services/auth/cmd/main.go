/**
 * @description
 * Main entry point for the Auth microservice.
 *
 * This file acts as the composition root for the application. It is responsible for:
 * - Loading configuration from environment variables.
 * - Establishing connections to external services (PostgreSQL, RabbitMQ).
 * - Wiring together all the application layers (repository, service, handlers, router).
 * - Starting the HTTP server to listen for requests.
 *
 * @dependencies
 * - "context"
 * - "log"
 * - "net/http"
 * - "os"
 * - "os/signal"
 * - "syscall"
 * - "time"
 * - All internal packages (api, app, config, store, pkg/rabbitmq)
 * - External libraries for pgxpool.
 */
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"transfa/services/auth/internal/api"
	"transfa/services/auth/internal/app"
	"transfa/services/auth/internal/config"
	"transfa/services/auth/internal/store"
	"transfa/services/auth/pkg/rabbitmq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Set the Clerk secret key
	clerk.SetKey(cfg.ClerkSecretKey)

	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection pool
	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}
	defer dbpool.Close()
	log.Println("Database connection pool established.")

	// Initialize RabbitMQ publisher
	publisher, err := rabbitmq.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("unable to create RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()
	log.Println("RabbitMQ publisher established.")

	// Wire application components
	repository := store.NewPostgresRepository(dbpool)
	service := app.NewService(repository, publisher, cfg)
	handler := api.NewAuthHandler(service)
	router := api.NewRouter(handler)

	// Set up and start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Auth Service is starting on port %s...", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Println("shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the requests it is currently handling
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
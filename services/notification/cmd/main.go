/**
 * @description
 * Main entry point for the Notification microservice.
 *
 * This file acts as the composition root for the application. It is responsible for:
 * - Loading configuration from environment variables.
 * - Establishing connections to external services (PostgreSQL, RabbitMQ).
 * - Wiring together all the application layers (repository, service, handlers, router).
 * - Starting the HTTP server to listen for requests, particularly webhooks.
 *
 * @dependencies
 * - Standard library packages for context, logging, HTTP, OS signals.
 * - External libraries for pgxpool and service-specific internal packages.
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

	"github.com/jackc/pgx/v5/pgxpool"
	"transfa/services/notification/internal/api"
	"transfa/services/notification/internal/app"
	"transfa/services/notification/internal/config"
	"transfa/services/notification/internal/store"
	"transfa/services/notification/pkg/rabbitmq"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

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
	handler := api.NewNotificationHandler(service)
	router := api.NewRouter(handler)

	// Set up and start HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("Notification Service is starting on port %s...", cfg.Port)
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
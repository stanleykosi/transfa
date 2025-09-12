/**
 * @description
 * Main entry point for the Account microservice.
 *
 * This file acts as the composition root for the application. It is responsible for:
 * - Loading configuration from environment variables.
 * - Establishing connections to external services (PostgreSQL).
 * - Initializing clients for other services (Anchor API).
 * - Wiring together all the application layers (repository, service, handlers).
 * - Starting the RabbitMQ consumer to process events asynchronously.
 * - Starting a minimal HTTP server for health checks.
 *
 * @dependencies
 * - Standard library packages for context, logging, HTTP, OS signals.
 * - External libraries for pgxpool, RabbitMQ, Viper.
 * - All internal packages for the account service.
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
	"transfa/services/account/internal/app"
	"transfa/services/account/internal/config"
	"transfa/services/account/internal/store"
	"transfa/services/account/pkg/anchor"
	"transfa/services/account/pkg/rabbitmq"
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

	// Wire application components
	repository := store.NewPostgresRepository(dbpool)
	anchorClient := anchor.NewClient(cfg.AnchorBaseURL, cfg.AnchorAPIKey)
	service := app.NewService(repository, anchorClient)

	// Initialize and start RabbitMQ consumer
	consumer, err := rabbitmq.NewConsumer(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to create RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	err = consumer.StartConsumer(
		ctx,
		cfg.CustomerVerifiedEx,
		cfg.CustomerVerifiedQueue,
		cfg.CustomerVerifiedRK,
		cfg.ConsumerTag,
		service.HandleCustomerVerifiedEvent,
	)
	if err != nil {
		log.Fatalf("failed to start RabbitMQ consumer: %v", err)
	}

	// Start a simple HTTP server for health checks in a goroutine
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		log.Printf("Account Service health check listening on port %s", cfg.Port)
		if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("health check server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully")

	// Allow some time for cleanup
	time.Sleep(2 * time.Second)
}
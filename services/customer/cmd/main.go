/**
 * @description
 * Main entry point for the Customer microservice.
 *
 * This file is the composition root. It's responsible for:
 * - Loading configuration from environment variables.
 * - Establishing connections to the PostgreSQL database and RabbitMQ.
 * - Initializing and wiring together all application components (repository, service, handlers, etc.).
 * - Starting the RabbitMQ consumer to listen for events.
 * - Starting a minimal HTTP server for health checks.
 *
 * @dependencies
 * - Standard library packages for context, logging, HTTP, OS signals.
 * - External libraries for pgxpool, RabbitMQ, and service-specific internal packages.
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
	"transfa/services/customer/internal/app"
	"transfa/services/customer/internal/config"
	"transfa/services/customer/internal/store"
	"transfa/services/customer/pkg/anchor"
	"transfa/services/customer/pkg/rabbitmq"
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
		cfg.UserCreatedEx,
		cfg.UserCreatedQueue,
		cfg.UserCreatedRK,
		cfg.ConsumerTag,
		service.HandleUserCreatedEvent,
	)
	if err != nil {
		log.Fatalf("failed to start RabbitMQ consumer: %v", err)
	}

	// Start a simple HTTP server for health checks in a goroutine
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		log.Printf("Customer Service health check listening on port %s", cfg.Port)
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
/**
 * @description
 * This file defines the interfaces (ports) for the application's core logic.
 * Following hexagonal architecture principles, these interfaces act as contracts
 * that external adapters (like the database repository or RabbitMQ publisher) must implement.
 * This decouples the core business logic from specific technologies.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "transfa/services/auth/internal/domain": Imports the core data models.
 */
package app

import (
	"context"
	"transfa/services/auth/internal/domain"
)

// Repository defines the interface for data persistence operations.
// Any database implementation must satisfy this interface.
type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

// Publisher defines the interface for publishing messages to a message broker.
// Any message broker implementation must satisfy this interface.
type Publisher interface {
	Publish(ctx context.Context, body []byte, exchange, routingKey string) error
	Close()
}
/**
 * @description
 * This file defines the interfaces (ports) for the application's core logic
 * in the Notification service. These interfaces act as contracts that external
 * adapters (like the database repository or RabbitMQ publisher) must implement.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "transfa/services/notification/internal/domain": Imports the core data models.
 */
package app

import (
	"context"

	"transfa/services/notification/internal/domain"
)

// Repository defines the interface for data persistence operations.
// It abstracts the database layer from the core application logic.
type Repository interface {
	GetUserByAnchorID(ctx context.Context, anchorID string) (*domain.User, error)
}

// Publisher defines the interface for publishing messages to a message broker.
type Publisher interface {
	Publish(ctx context.Context, body []byte, exchange, routingKey string) error
	Close()
}
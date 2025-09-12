/**
 * @description
 * This file defines the interfaces (ports) for the Customer service's application logic.
 * These interfaces define the contracts for external dependencies, such as the database
 * and the Anchor API client, allowing for a clean separation of concerns and easier testing.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "github.com/google/uuid": For user identifiers.
 * - "transfa/services/customer/internal/domain": For event data structures.
 */
package app

import (
	"context"

	"github.com/google/uuid"
	"transfa/services/customer/internal/domain"
)

// Repository defines the interface for data persistence operations.
type Repository interface {
	UpdateUserWithAnchorID(ctx context.Context, userID uuid.UUID, anchorCustomerID string) error
}

// AnchorClient defines the interface for communicating with the Anchor BaaS API.
type AnchorClient interface {
	CreateIndividualCustomer(ctx context.Context, event domain.UserCreatedEvent) (string, error)
	TriggerIndividualVerification(ctx context.Context, anchorCustomerID string, kycDetails *domain.KYCDetails) error
	// CreateBusinessCustomer and TriggerBusinessVerification would be defined here as well.
}
/**
 * @description
 * This file defines the interfaces (ports) for the Account service's application logic.
 * These interfaces define the contracts for external dependencies, such as the database
 * and the Anchor API client, allowing for a clean separation of concerns and easier testing.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "github.com/google/uuid": For user identifiers.
 * - "transfa/services/account/internal/domain": For core data models.
 */
package app

import (
	"context"

	"github.com/google/uuid"
	"transfa/services/account/internal/domain"
)

// Repository defines the interface for data persistence operations.
type Repository interface {
	CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

// AnchorClient defines the interface for communicating with the Anchor BaaS API.
type AnchorClient interface {
	CreateDepositAccount(ctx context.Context, anchorCustomerID, customerType, productName string) (string, error)
}
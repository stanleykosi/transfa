/**
 * @description
 * This file contains the core business logic for the Account service. The Service struct
 * orchestrates the process of handling a `customer.verified` event by interacting with the
 * Anchor client and the database repository to create a user wallet.
 *
 * @dependencies
 * - Go standard libraries: "context", "encoding/json", "fmt", "log"
 * - "github.com/rabbitmq/amqp091-go": For message handling.
 * - "transfa/services/account/internal/domain": For core data models and events.
 */
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"transfa/services/account/internal/domain"
)

// Service provides the application's business logic for account management.
type Service struct {
	repo         Repository
	anchorClient AnchorClient
}

// NewService creates a new application service.
func NewService(repo Repository, anchorClient AnchorClient) *Service {
	return &Service{
		repo:         repo,
		anchorClient: anchorClient,
	}
}

// HandleCustomerVerifiedEvent is the message handler for `customer.verified` events.
// It orchestrates creating the Anchor DepositAccount and saving it to the local DB.
func (s *Service) HandleCustomerVerifiedEvent(ctx context.Context, msg amqp091.Delivery) error {
	var event domain.CustomerVerifiedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal CustomerVerifiedEvent: %w", err)
	}

	log.Printf("Processing CustomerVerifiedEvent for UserID: %s", event.UserID)

	// Step 1: Fetch the user from our DB to get their account type.
	user, err := s.repo.GetUserByID(ctx, event.UserID)
	if err != nil {
		// If the user doesn't exist, we can't proceed. This would indicate an issue upstream.
		return fmt.Errorf("failed to get user by ID %s: %w", event.UserID, err)
	}

	// Step 2: Determine the correct product and customer types based on the user's account type.
	var productName, customerType string
	switch user.AccountType {
	case "personal":
		productName = "SAVINGS"
		customerType = "IndividualCustomer"
	case "merchant":
		productName = "CURRENT"
		customerType = "BusinessCustomer"
	default:
		return fmt.Errorf("unknown account type '%s' for user %s", user.AccountType, user.ID)
	}

	// Step 3: Call the Anchor API to create the DepositAccount.
	anchorAccountID, err := s.anchorClient.CreateDepositAccount(ctx, event.AnchorCustomerID, customerType, productName)
	if err != nil {
		return fmt.Errorf("failed to create deposit account in anchor for user %s: %w", event.UserID, err)
	}

	log.Printf("Successfully created Anchor DepositAccount with ID: %s for UserID: %s", anchorAccountID, event.UserID)

	// Step 4: Create the account record in our local database.
	newAccount := &domain.Account{
		UserID:          user.ID,
		AnchorAccountID: anchorAccountID,
		AccountPurpose:  "main_wallet",
		Status:          "active",
		Balance:         0,
	}

	if _, err := s.repo.CreateAccount(ctx, newAccount); err != nil {
		// This is a critical error. We have an orphaned Anchor account.
		// Requires robust retry logic or manual intervention.
		return fmt.Errorf("CRITICAL: failed to save created account for user %s with anchor_account_id %s: %w", user.ID, anchorAccountID, err)
	}

	log.Printf("Successfully stored new account record for user %s", user.ID)
	// The technical spec mentions a final `account.opened` webhook from Anchor.
	// The Notification service would listen for this and could send a final "Welcome!" push notification.

	return nil
}
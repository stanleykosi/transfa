/**
 * @description
 * This file contains the core business logic for the Customer service. The Service struct
 * orchestrates the process of handling a `user.created` event by interacting with the
 * Anchor client and the database repository.
 *
 * @dependencies
 * - "context", "encoding/json", "fmt", "log"
 * - "github.com/rabbitmq/amqp091-go": For message handling.
 * - "transfa/services/customer/internal/domain": For core data models and events.
 */
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"transfa/services/customer/internal/domain"
)

// Service provides the application's business logic.
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

// HandleUserCreatedEvent is the message handler for `user.created` events.
// It orchestrates creating the customer in Anchor, updating the local DB, and triggering KYC.
func (s *Service) HandleUserCreatedEvent(ctx context.Context, msg amqp091.Delivery) error {
	var event domain.UserCreatedEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal UserCreatedEvent: %w", err)
	}

	log.Printf("Processing UserCreatedEvent for UserID: %s", event.UserID)

	// Step 1: Create customer in Anchor based on account type.
	var anchorCustomerID string
	var err error

	switch event.AccountType {
	case "personal":
		if event.KYCDetails == nil {
			return fmt.Errorf("KYCDetails are required for personal account type")
		}
		anchorCustomerID, err = s.anchorClient.CreateIndividualCustomer(ctx, event)
		if err != nil {
			return fmt.Errorf("failed to create individual customer in anchor: %w", err)
		}
	case "merchant":
		// NOTE: As per the spec, the event only contains business_name and rc_number,
		// which is insufficient for Anchor's detailed business customer creation.
		// A placeholder or a call to a simplified endpoint would be needed here.
		// For now, we log this as a work-in-progress.
		log.Printf("Merchant onboarding for UserID %s is not fully implemented due to insufficient data.", event.UserID)
		// anchorCustomerID, err = s.anchorClient.CreateBusinessCustomer(ctx, event)
		return nil // Acknowledge message to prevent requeue loop
	default:
		return fmt.Errorf("unknown account type: %s", event.AccountType)
	}

	log.Printf("Successfully created Anchor customer with ID: %s for UserID: %s", anchorCustomerID, event.UserID)

	// Step 2: Update our local user record with the Anchor Customer ID.
	if err := s.repo.UpdateUserWithAnchorID(ctx, event.UserID, anchorCustomerID); err != nil {
		// This is a critical error. If this fails, we have an orphaned Anchor customer.
		// This might require a retry mechanism or manual intervention.
		return fmt.Errorf("CRITICAL: failed to update user %s with anchor id %s: %w", event.UserID, anchorCustomerID, err)
	}
	log.Printf("Successfully updated user %s with anchor_customer_id", event.UserID)

	// Step 3: Trigger the verification process in Anchor.
	if event.AccountType == "personal" {
		if err := s.anchorClient.TriggerIndividualVerification(ctx, anchorCustomerID, event.KYCDetails); err != nil {
			// This is less critical; it can be retried. We still acknowledge the message.
			log.Printf("WARNING: Failed to trigger KYC verification for anchor customer %s: %v", anchorCustomerID, err)
		}
	}
	// Similar logic for business verification would go here.

	return nil
}
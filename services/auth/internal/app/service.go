/**
 * @description
 * This file contains the core business logic for the Auth service.
 * The Service struct orchestrates operations by coordinating between the domain models,
 * the repository (for data persistence), and the publisher (for event messaging).
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "encoding/json": For serializing event data.
 * - "log": For logging information and errors.
 * - "github.com/google/uuid": To generate UUIDs for new users.
 * - "transfa/services/auth/internal/config": Imports app configuration.
 * - "transfa/services/auth/internal/domain": Imports the core data models.
 */
package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/google/uuid"
	"transfa/services/auth/internal/config"
	"transfa/services/auth/internal/domain"
)

// Service provides the application's business logic.
type Service struct {
	repo      Repository
	publisher Publisher
	config    config.Config
}

// NewService creates a new application service.
func NewService(repo Repository, publisher Publisher, cfg config.Config) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
		config:    cfg,
	}
}

// OnboardUser handles the business logic for creating a new user.
// It creates a user record in the database and publishes a `user.created` event.
func (s *Service) OnboardUser(ctx context.Context, clerkID string, req domain.OnboardingRequest) (*domain.User, error) {
	// 1. Construct the user object from the request.
	newUser := &domain.User{
		ID:          uuid.New(), // Generate a new UUID for the user.
		ClerkID:     clerkID,
		Username:    req.Username,
		AccountType: req.AccountType,
	}

	// Set the allow_sending flag based on account type.
	if req.AccountType == "merchant" {
		newUser.AllowSending = false
	} else {
		newUser.AllowSending = true
	}

	// 2. Persist the user to the database.
	createdUser, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		log.Printf("Failed to create user in repository: %v", err)
		return nil, err
	}

	// 3. Prepare and publish the `user.created` event.
	event := domain.UserCreatedEvent{
		UserID:      createdUser.ID,
		ClerkID:     createdUser.ClerkID,
		AccountType: createdUser.AccountType,
		KYCDetails:  req.KYCDetails,
		KYBDetails:  req.KYBDetails,
	}

	eventBody, err := json.Marshal(event)
	if err != nil {
		// Log the error but don't fail the whole operation, as the user is already created.
		// This might be a scenario for a dead-letter queue or another recovery mechanism.
		log.Printf("ERROR: Failed to marshal UserCreatedEvent for user %s: %v", createdUser.ID, err)
		return createdUser, nil
	}

	err = s.publisher.Publish(ctx, eventBody, s.config.UserCreatedEx, s.config.UserCreatedRK)
	if err != nil {
		// Same as above, log the error but the user creation itself was successful.
		log.Printf("ERROR: Failed to publish UserCreatedEvent for user %s: %v", createdUser.ID, err)
	} else {
		log.Printf("Successfully published UserCreatedEvent for user %s", createdUser.ID)
	}


	return createdUser, nil
}
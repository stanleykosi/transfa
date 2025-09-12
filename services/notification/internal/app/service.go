/**
 * @description
 * This file contains the core business logic for the Notification service.
 * The Service struct orchestrates the processing of incoming webhooks, including
 * signature verification, event parsing, and publishing new internal events.
 *
 * @dependencies
 * - Go standard libraries: "context", "crypto/hmac", "crypto/sha1", "crypto/subtle", "encoding/base64", "encoding/json", "fmt", "log"
 * - Internal packages: "config", "domain", "store" for application-specific logic and models.
 */
package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"transfa/services/notification/internal/config"
	"transfa/services/notification/internal/domain"
	"transfa/services/notification/internal/store"
)

// Service provides the application's business logic for notifications and webhooks.
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

// ProcessAnchorWebhook verifies, parses, and acts upon an incoming webhook from Anchor.
func (s *Service) ProcessAnchorWebhook(ctx context.Context, payload []byte, signature string) error {
	// 1. Verify the webhook signature to ensure it's from Anchor.
	if err := s.verifySignature(payload, signature); err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	// 2. Unmarshal the payload into our defined struct.
	var webhook domain.AnchorWebhookPayload
	if err := json.Unmarshal(payload, &webhook); err != nil {
		return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}

	log.Printf("Received valid Anchor webhook of type: %s", webhook.Data.Type)

	// 3. Process the event based on its type.
	switch webhook.Data.Type {
	case "customer.identification.approved":
		return s.handleCustomerIdentificationApproved(ctx, webhook)
	case "customer.identification.rejected":
		return s.handleCustomerIdentificationRejected(ctx, webhook)
	default:
		log.Printf("Unhandled Anchor event type: %s", webhook.Data.Type)
		return nil // Acknowledge unhandled events to prevent requeues.
	}
}

// verifySignature calculates the HMAC-SHA1 signature of the payload and compares it securely
// with the signature provided in the header.
func (s *Service) verifySignature(payload []byte, signatureHeader string) error {
	if s.config.AnchorWebhookSecret == "" {
		return errors.New("anchor webhook secret is not configured")
	}

	mac := hmac.New(sha1.New, []byte(s.config.AnchorWebhookSecret))
	mac.Write(payload)
	expectedMAC := mac.Sum(nil)
	expectedSignature := base64.StdEncoding.EncodeToString(expectedMAC)

	// Use subtle.ConstantTimeCompare to prevent timing attacks.
	if subtle.ConstantTimeCompare([]byte(signatureHeader), []byte(expectedSignature)) != 1 {
		return errors.New("signatures do not match")
	}

	return nil
}

// handleCustomerIdentificationApproved processes a successful KYC/KYB webhook.
func (s *Service) handleCustomerIdentificationApproved(ctx context.Context, webhook domain.AnchorWebhookPayload) error {
	anchorCustomerID := webhook.Data.Relationships.Customer.ID
	if anchorCustomerID == "" {
		return errors.New("missing anchor_customer_id in webhook payload")
	}

	// We need to find our internal user ID from the anchor customer ID.
	user, err := s.repo.GetUserByAnchorID(ctx, anchorCustomerID)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			log.Printf("WARNING: Received approval for an unknown Anchor customer: %s", anchorCustomerID)
			return nil // Acknowledge to prevent requeue, but log as a warning.
		}
		return fmt.Errorf("failed to get user by anchor ID: %w", err)
	}

	// Create and publish the internal event.
	event := domain.CustomerVerifiedEvent{
		UserID:           user.ID,
		AnchorCustomerID: anchorCustomerID,
	}

	eventBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal CustomerVerifiedEvent: %w", err)
	}

	err = s.publisher.Publish(ctx, eventBody, s.config.CustomerVerifiedEx, s.config.CustomerVerifiedRK)
	if err != nil {
		return fmt.Errorf("failed to publish CustomerVerifiedEvent: %w", err)
	}

	log.Printf("Successfully published CustomerVerifiedEvent for UserID: %s", user.ID)
	return nil
}

// handleCustomerIdentificationRejected processes a failed KYC/KYB webhook.
func (s *Service) handleCustomerIdentificationRejected(ctx context.Context, webhook domain.AnchorWebhookPayload) error {
	anchorCustomerID := webhook.Data.Relationships.Customer.ID
	if anchorCustomerID == "" {
		return errors.New("missing anchor_customer_id in rejected webhook payload")
	}

	log.Printf("Received KYC rejection for Anchor Customer ID: %s", anchorCustomerID)

	user, err := s.repo.GetUserByAnchorID(ctx, anchorCustomerID)
	if err != nil {
		log.Printf("WARNING: Received rejection for an unknown Anchor customer: %s", anchorCustomerID)
		return nil
	}

	// TODO: Parse the rejection reason from `webhook.Data.Attributes`.
	// For now, we'll use a generic reason.
	rejectionReason := "KYC details could not be verified."

	event := domain.CustomerVerificationRejectedEvent{
		UserID:           user.ID,
		AnchorCustomerID: anchorCustomerID,
		Reason:           rejectionReason,
	}
	eventBody, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal CustomerVerificationRejectedEvent: %w", err)
	}

	err = s.publisher.Publish(ctx, eventBody, s.config.CustomerVerificationRejectedEx, s.config.CustomerVerificationRejectedRK)
	if err != nil {
		return fmt.Errorf("failed to publish CustomerVerificationRejectedEvent: %w", err)
	}

	log.Printf("Successfully published CustomerVerificationRejectedEvent for UserID: %s", user.ID)
	// TODO: In a future step, consume this event to send a push notification to the user.
	return nil
}
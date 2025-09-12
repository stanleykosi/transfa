/**
 * @description
 * This file defines the domain models and Data Transfer Objects (DTOs) for the
 * Notification service. It includes structures for parsing incoming webhooks from Anchor
 * and for creating outgoing events to be published to RabbitMQ.
 *
 * @dependencies
 * - "time": For timestamping.
 * - "github.com/google/uuid": For universally unique identifiers.
 */
package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AnchorRelationshipData represents the content of a relationship object in a JSON:API payload.
type AnchorRelationshipData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// AnchorRelationships defines the relationships block in an Anchor webhook.
type AnchorRelationships struct {
	Customer AnchorRelationshipData `json:"customer"`
}

// AnchorWebhookData is the main "data" object within an Anchor webhook payload.
type AnchorWebhookData struct {
	ID            string              `json:"id"`
	Type          string              `json:"type"` // e.g., "customer.identification.approved"
	Attributes    json.RawMessage     `json:"attributes"`
	Relationships AnchorRelationships `json:"relationships"`
}

// AnchorWebhookPayload is the top-level structure for an incoming webhook from Anchor.
type AnchorWebhookPayload struct {
	Data AnchorWebhookData `json:"data"`
}

// CustomerVerifiedEvent is the payload published to RabbitMQ when a customer's KYC is approved.
type CustomerVerifiedEvent struct {
	UserID           uuid.UUID `json:"user_id"`
	AnchorCustomerID string    `json:"anchor_customer_id"`
}

// CustomerVerificationRejectedEvent is the payload for when KYC fails.
type CustomerVerificationRejectedEvent struct {
	UserID           uuid.UUID `json:"user_id"`
	AnchorCustomerID string    `json:"anchor_customer_id"`
	Reason           string    `json:"reason"` // We can add more details from the webhook later
}

// User is a simplified representation of our user table, needed to find the
// internal user ID from an Anchor customer ID.
type User struct {
	ID uuid.UUID `db:"id"`
}
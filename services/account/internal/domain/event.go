/**
 * @description
 * This file defines the structure of events that the Account service consumes.
 * Specifically, it defines the `CustomerVerifiedEvent`, which is the message payload
 * received from the Notification service via RabbitMQ when a user's KYC/KYB is approved.
 *
 * @dependencies
 * - "github.com/google/uuid": Used for universally unique identifiers.
 */
package domain

import "github.com/google/uuid"

// CustomerVerifiedEvent is the message structure for the `customer.verified` event.
// It contains the necessary information for the Account service to create
// a new wallet (DepositAccount) for the user in the Anchor BaaS.
type CustomerVerifiedEvent struct {
	UserID           uuid.UUID `json:"user_id"`
	AnchorCustomerID string    `json:"anchor_customer_id"`
}
/**
 * @description
 * This file defines the structure of events that the Customer service consumes.
 * Specifically, it defines the `UserCreatedEvent`, which is the message payload
 * received from the Auth service via RabbitMQ when a new user completes onboarding.
 *
 * @dependencies
 * - "github.com/google/uuid": Used for universally unique identifiers.
 */
package domain

import "github.com/google/uuid"

// KYCDetails holds the Know Your Customer information for personal users.
// This structure must match the one published by the Auth service.
type KYCDetails struct {
	FullName    string `json:"full_name"`
	BVN         string `json:"bvn"`
	DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

// KYBDetails holds the Know Your Business information for merchant users.
// This structure must match the one published by the Auth service.
type KYBDetails struct {
	BusinessName string `json:"business_name"`
	RCNumber     string `json:"rc_number"`
}

// UserCreatedEvent is the message structure for the `user.created` event.
// It contains all the necessary information for the Customer service to create
// a customer record in the Anchor BaaS.
type UserCreatedEvent struct {
	UserID      uuid.UUID   `json:"user_id"`
	ClerkID     string      `json:"clerk_id"`
	AccountType string      `json:"account_type"`
	KYCDetails  *KYCDetails `json:"kyc_details,omitempty"`
	KYBDetails  *KYBDetails `json:"kyb_details,omitempty"`
}
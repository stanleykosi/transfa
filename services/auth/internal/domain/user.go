/**
 * @description
 * This file defines the core domain models and Data Transfer Objects (DTOs)
 * related to user onboarding in the Auth service.
 *
 * Key features:
 * - `OnboardingRequest`: Defines the JSON structure for the POST /onboarding endpoint.
 * - `KYCDetails` & `KYBDetails`: Specific structures for personal and merchant identity information.
 * - `User`: Represents the user entity as it's stored in the database.
 * - `UserCreatedEvent`: Defines the structure of the event published to RabbitMQ after user creation.
 *
 * @dependencies
 * - "time": Used for timestamping records.
 * - "github.com/google/uuid": Used for universally unique identifiers.
 */
package domain

import (
	"time"
	"github.com/google/uuid"
)

// KYCDetails holds the Know Your Customer information for personal users.
type KYCDetails struct {
	FullName    string `json:"full_name"`
	BVN         string `json:"bvn"`
	DateOfBirth string `json:"date_of_birth"` // YYYY-MM-DD
}

// KYBDetails holds the Know Your Business information for merchant users.
type KYBDetails struct {
	BusinessName string `json:"business_name"`
	RCNumber     string `json:"rc_number"`
}

// OnboardingRequest is the expected JSON body for the user onboarding endpoint.
type OnboardingRequest struct {
	Username    string      `json:"username"`
	AccountType string      `json:"account_type"` // "personal" or "merchant"
	KYCDetails  *KYCDetails `json:"kyc_details,omitempty"`
	KYBDetails  *KYBDetails `json:"kyb_details,omitempty"`
}

// User represents the core user profile in the Transfa system.
// It maps directly to the `users` table in the database.
type User struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ClerkID     string    `json:"clerk_id" db:"clerk_id"`
	Username    string    `json:"username" db:"username"`
	AccountType string    `json:"account_type" db:"account_type"`
	AllowSending bool     `json:"allow_sending" db:"allow_sending"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	// The following fields will be populated by other services.
	// AnchorCustomerID string    `json:"anchor_customer_id" db:"anchor_customer_id"`
	// KYCStatus        string    `json:"kyc_status" db:"kyc_status"`
	// ProfileImageURL  *string   `json:"profile_image_url,omitempty" db:"profile_image_url"`
}

// UserCreatedEvent represents the payload published to RabbitMQ after a user is created.
// This event triggers the customer creation flow in the Customer service.
type UserCreatedEvent struct {
	UserID      uuid.UUID    `json:"user_id"`
	ClerkID     string       `json:"clerk_id"`
	AccountType string       `json:"account_type"`
	KYCDetails  *KYCDetails  `json:"kyc_details,omitempty"`
	KYBDetails  *KYBDetails  `json:"kyb_details,omitempty"`
}
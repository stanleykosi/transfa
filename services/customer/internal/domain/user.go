/**
 * @description
 * This file defines the core domain models for user-related entities within the Customer service.
 * These structs represent the application's core data structures and are used across different
 * layers (API, application logic, and storage).
 *
 * Key features:
 * - `User`: Represents a user profile, linking authentication (Clerk), BaaS (Anchor), and app-specific data.
 * - `UserSettings`: Stores user-specific preferences, such as default beneficiaries.
 *
 * @dependencies
 * - "time": Used for timestamping records.
 * - "github.com/google/uuid": Used for universally unique identifiers as primary keys.
 */

package domain

import (
"time"

"github.com/google/uuid"
)

// User represents the core user profile in the Transfa system.
// It maps directly to the `users` table in the database.
type User struct {
ID               uuid.UUID `json:"id" db:"id"`
ClerkID          string    `json:"clerk_id" db:"clerk_id"`
Username         string    `json:"username" db:"username"`
AccountType      string    `json:"account_type" db:"account_type"`
AnchorCustomerID string    `json:"anchor_customer_id" db:"anchor_customer_id"`
KYCStatus        string    `json:"kyc_status" db:"kyc_status"`
ProfileImageURL  *string   `json:"profile_image_url,omitempty" db:"profile_image_url"`
AllowSending     bool      `json:"allow_sending" db:"allow_sending"`
CreatedAt        time.Time `json:"created_at" db:"created_at"`
UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// UserSettings represents user-specific preferences.
// It maps directly to the `user_settings` table in the database.
type UserSettings struct {
UserID               uuid.UUID  `json:"user_id" db:"user_id"`
DefaultBeneficiaryID *uuid.UUID `json:"default_beneficiary_id,omitempty" db:"default_beneficiary_id"`
UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
}
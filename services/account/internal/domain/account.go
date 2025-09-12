/**
 * @description
 * This file defines the core domain model for an Account within the Account service.
 * An Account represents a user's wallet or other special-purpose accounts (e.g., for Money Drops)
 * within the Transfa system. This model corresponds to an Anchor `DepositAccount`.
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

// Account represents a user's wallet within the Transfa application.
// This could be their main wallet or a special-purpose wallet like for a Money Drop.
// It maps directly to the `accounts` table in the database.
type Account struct {
ID              uuid.UUID `json:"id" db:"id"`
UserID          uuid.UUID `json:"user_id" db:"user_id"`
AnchorAccountID string    `json:"anchor_account_id" db:"anchor_account_id"`
AccountPurpose  string    `json:"account_purpose" db:"account_purpose"`
Balance         int64     `json:"balance" db:"balance"` // Stored in kobo
Status          string    `json:"status" db:"status"`
CreatedAt       time.Time `json:"created_at" db:"created_at"`
UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
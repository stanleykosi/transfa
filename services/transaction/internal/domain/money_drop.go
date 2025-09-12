/**
 * @description
 * This file defines the domain models related to the Money Drop feature within the Transaction service.
 *
 * Key features:
 * - `MoneyDrop`: Represents a Money Drop event, including its funding, rules, and current state.
 * - `MoneyDropClaim`: Records an individual user's claim against a Money Drop to ensure the
 *   "one claim per person" rule is enforced.
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

// MoneyDrop represents a created Money Drop instance, containing its rules and current status.
// It maps directly to the `money_drops` table in the database.
type MoneyDrop struct {
ID                  uuid.UUID `json:"id" db:"id"`
CreatorUserID       uuid.UUID `json:"creator_user_id" db:"creator_user_id"`
FundingAccountID    uuid.UUID `json:"funding_account_id" db:"funding_account_id"`
TotalAmount         int64     `json:"total_amount" db:"total_amount"`                 // Stored in kobo
AmountPerClaim      int64     `json:"amount_per_claim" db:"amount_per_claim"`         // Stored in kobo
TotalClaimsAllowed  int       `json:"total_claims_allowed" db:"total_claims_allowed"`
ClaimsMadeCount     int       `json:"claims_made_count" db:"claims_made_count"`
Status              string    `json:"status" db:"status"`
ExpiryTimestamp     time.Time `json:"expiry_timestamp" db:"expiry_timestamp"`
CreatedAt           time.Time `json:"created_at" db:"created_at"`
UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// MoneyDropClaim represents a record of a user claiming a share from a Money Drop.
// This is used to prevent duplicate claims.
// It maps directly to the `money_drop_claims` table.
type MoneyDropClaim struct {
ID             uuid.UUID `json:"id" db:"id"`
MoneyDropID    uuid.UUID `json:"money_drop_id" db:"money_drop_id"`
ClaimantUserID uuid.UUID `json:"claimant_user_id" db:"claimant_user_id"`
ClaimedAt      time.Time `json:"claimed_at" db:"claimed_at"`
}
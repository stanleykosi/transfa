/**
 * @description
 * This file defines the domain model for a Payment Request within the Transaction service.
 * A Payment Request is a user-generated request for a specific amount of money, which can
 * be fulfilled by another user.
 *
 * @dependencies
 * - "database/sql": Used for nullable time types.
 * - "time": Used for timestamping records.
 * - "github.com/google/uuid": Used for universally unique identifiers and nullable UUIDs.
 */
package domain

import (
"database/sql"
"time"

"github.com/google/uuid"
)

// PaymentRequest represents a user-created request for payment.
// It maps directly to the `payment_requests` table in the database.
type PaymentRequest struct {
ID                      uuid.UUID     `json:"id" db:"id"`
CreatorUserID           uuid.UUID     `json:"creator_user_id" db:"creator_user_id"`
Amount                  int64         `json:"amount" db:"amount"` // Stored in kobo
Description             *string       `json:"description,omitempty" db:"description"`
ImageURL                *string       `json:"image_url,omitempty" db:"image_url"`
Status                  string        `json:"status" db:"status"`
FulfilledAt             sql.NullTime  `json:"fulfilled_at,omitempty" db:"fulfilled_at"`
FulfilledByTransactionID uuid.NullUUID `json:"fulfilled_by_transaction_id,omitempty" db:"fulfilled_by_transaction_id"`
CreatedAt               time.Time     `json:"created_at" db:"created_at"`
UpdatedAt               time.Time     `json:"updated_at" db:"updated_at"`
}
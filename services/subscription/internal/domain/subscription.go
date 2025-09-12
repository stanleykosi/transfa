/**
 * @description
 * This file defines the domain model for a Subscription within the Subscription service.
 * This model is crucial for managing user entitlements, transfer routing logic, and recurring billing.
 *
 * @dependencies
 * - "database/sql": Used for nullable time types.
 * - "time": Used for timestamping records.
 * - "github.com/google/uuid": Used for universally unique identifiers as primary keys.
 */

package domain

import (
"database/sql"
"time"

"github.com/google/uuid"
)

// Subscription represents a user's subscription status and entitlements.
// It maps directly to the `subscriptions` table in the database.
type Subscription struct {
ID                           uuid.UUID    `json:"id" db:"id"`
UserID                       uuid.UUID    `json:"user_id" db:"user_id"`
Status                       string       `json:"status" db:"status"`
AutoRenew                    bool         `json:"auto_renew" db:"auto_renew"`
CurrentPeriodEndsAt          sql.NullTime `json:"current_period_ends_at,omitempty" db:"current_period_ends_at"`
MonthlyExternalTransfersUsed int          `json:"monthly_external_transfers_used" db:"monthly_external_transfers_used"`
CreatedAt                    time.Time    `json:"created_at" db:"created_at"`
UpdatedAt                    time.Time    `json:"updated_at" db:"updated_at"`
}
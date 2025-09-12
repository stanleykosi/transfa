/**
 * @description
 * This file defines the domain model for a Beneficiary within the Customer service.
 * A Beneficiary represents an external bank account that a user has saved for self-transfers
 * (withdrawals). This model maps directly to an Anchor `CounterParty` resource.
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

// Beneficiary represents a user's saved external bank account for withdrawals.
// It maps directly to the `beneficiaries` table in the database.
type Beneficiary struct {
ID                  uuid.UUID `json:"id" db:"id"`
UserID              uuid.UUID `json:"user_id" db:"user_id"`
AnchorCounterpartyID string    `json:"anchor_counterparty_id" db:"anchor_counterparty_id"`
AccountName         string    `json:"account_name" db:"account_name"`
AccountNumber       string    `json:"account_number" db:"account_number"`
BankName            string    `json:"bank_name" db:"bank_name"`
BankCode            string    `json:"bank_code" db:"bank_code"`
IsDefault           bool      `json:"is_default" db:"is_default"`
CreatedAt           time.Time `json:"created_at" db:"created_at"`
UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
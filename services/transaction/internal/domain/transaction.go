/**
 * @description
 * This file defines the core domain model for a Transaction within the Transaction service.
 * The Transaction struct is a comprehensive record of any money movement within the
 * Transfa application, covering everything from P2P payments to wallet funding.
 *
 * @dependencies
 * - "time": Used for timestamping records.
 * - "github.com/google/uuid": Used for universally unique identifiers and nullable UUIDs.
 */
package domain

import (
"time"

"github.com/google/uuid"
)

// Transaction represents a single financial transaction in the system.
// It maps directly to the `transactions` table in the database.
// Note the use of `uuid.NullUUID` for optional foreign keys, accommodating
// different transaction types (e.g., wallet funding has no sender).
type Transaction struct {
ID                       uuid.UUID     `json:"id" db:"id"`
SenderUserID             uuid.NullUUID `json:"sender_user_id,omitempty" db:"sender_user_id"`
RecipientUserID          uuid.NullUUID `json:"recipient_user_id,omitempty" db:"recipient_user_id"`
SourceAccountID          uuid.NullUUID `json:"source_account_id,omitempty" db:"source_account_id"`
DestinationAccountID     uuid.NullUUID `json:"destination_account_id,omitempty" db:"destination_account_id"`
DestinationBeneficiaryID uuid.NullUUID `json:"destination_beneficiary_id,omitempty" db:"destination_beneficiary_id"`
AnchorTransferID         *string       `json:"anchor_transfer_id,omitempty" db:"anchor_transfer_id"`
Type                     string        `json:"type" db:"type"`
Amount                   int64         `json:"amount" db:"amount"` // Stored in kobo
Fee                      int64         `json:"fee" db:"fee"`       // Stored in kobo
Status                   string        `json:"status" db:"status"`
Description              *string       `json:"description,omitempty" db:"description"`
Category                 *string       `json:"category,omitempty" db:"category"`
CreatedAt                time.Time     `json:"created_at" db:"created_at"`
UpdatedAt                time.Time     `json:"updated_at" db:"updated_at"`
}
/**
 * @description
 * This file defines a simplified User domain model for the Account service.
 * This service only needs to know a user's ID and account type to make decisions
 * about what kind of wallet to create.
 *
 * @dependencies
 * - "github.com/google/uuid": Used for the user's unique identifier.
 */
package domain

import "github.com/google/uuid"

// User represents a minimal user profile needed by the Account service.
// It is used to fetch the account_type from the database.
type User struct {
	ID          uuid.UUID `db:"id"`
	AccountType string    `db:"account_type"`
}
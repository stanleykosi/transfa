/**
 * @description
 * This file provides the PostgreSQL implementation of the Repository interface for the Customer service.
 * It handles all direct database interactions, such as updating user records.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "fmt": For formatting error messages.
 * - "github.com/google/uuid": For user identifiers.
 * - "github.com/jackc/pgx/v5/pgxpool": The PostgreSQL driver and connection pool.
 */
package store

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository is the concrete implementation for database operations.
type PostgresRepository struct {
	db *pgxpool.Pool
}

// NewPostgresRepository creates a new instance of the repository.
func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// UpdateUserWithAnchorID updates a user's record with their new Anchor Customer ID.
// This is a critical step after successfully creating the customer in the BaaS.
func (r *PostgresRepository) UpdateUserWithAnchorID(ctx context.Context, userID uuid.UUID, anchorCustomerID string) error {
	query := `
        UPDATE public.users
        SET anchor_customer_id = $1
        WHERE id = $2
    `
	cmdTag, err := r.db.Exec(ctx, query, anchorCustomerID, userID)
	if err != nil {
		return fmt.Errorf("failed to execute update user query: %w", err)
	}

	if cmdTag.RowsAffected() != 1 {
		return fmt.Errorf("expected 1 row to be affected, but got %d", cmdTag.RowsAffected())
	}

	return nil
}
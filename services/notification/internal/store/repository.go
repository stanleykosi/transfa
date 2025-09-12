/**
 * @description
 * This file provides the PostgreSQL implementation of the Repository interface for the
 * Notification service. It encapsulates all database-specific logic, primarily for
 * retrieving user information based on external identifiers.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "errors": For handling specific database errors like "no rows".
 * - "fmt": For formatting error messages.
 * - "github.com/jackc/pgx/v5": For checking specific database errors.
 * - "github.com/jackc/pgx/v5/pgxpool": For managing the database connection pool.
 * - "transfa/services/notification/internal/domain": Imports the User model.
 */
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"transfa/services/notification/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

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

// GetUserByAnchorID retrieves a user's internal ID from the database using their
// unique Anchor Customer ID.
func (r *PostgresRepository) GetUserByAnchorID(ctx context.Context, anchorID string) (*domain.User, error) {
	query := `SELECT id FROM public.users WHERE anchor_customer_id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, anchorID).Scan(&user.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: with anchor_customer_id %s", ErrUserNotFound, anchorID)
		}
		return nil, fmt.Errorf("failed to query user by anchor id: %w", err)
	}

	return &user, nil
}
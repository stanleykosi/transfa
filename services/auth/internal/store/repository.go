/**
 * @description
 * This file provides the PostgreSQL implementation of the Repository interface.
 * It encapsulates all database-specific logic for the Auth service, using the `pgx`
 * driver for efficient and safe interaction with the PostgreSQL database.
 *
 * @dependencies
 * - "context": For passing request-scoped data and cancellation signals.
 * - "github.com/jackc/pgx/v5/pgxpool": For managing the database connection pool.
 * - "transfa/services/auth/internal/domain": Imports the core User model.
 */
package store

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"transfa/services/auth/internal/domain"
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

// CreateUser inserts a new user record into the public.users table.
// It returns the newly created user with fields populated by the database (like id, created_at).
func (r *PostgresRepository) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
        INSERT INTO public.users (id, clerk_id, username, account_type, allow_sending)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, clerk_id, username, account_type, allow_sending, created_at, updated_at
    `

	// The user ID is linked to auth.users and should be the same UUID.
	// For this operation, we assume the trigger/RLS will handle auth linking.
	// The Clerk User ID from the JWT 'sub' claim will be the 'clerk_id'.

	var createdUser domain.User
	err := r.db.QueryRow(ctx, query,
		user.ID,
		user.ClerkID,
		user.Username,
		user.AccountType,
		user.AllowSending,
	).Scan(
		&createdUser.ID,
		&createdUser.ClerkID,
		&createdUser.Username,
		&createdUser.AccountType,
		&createdUser.AllowSending,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error creating user in database: %v", err)
		// Here you could check for specific pgx errors, e.g., unique constraint violation
		// and return a more specific application-level error.
		return nil, err
	}

	return &createdUser, nil
}
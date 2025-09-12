/**
 * @description
 * This file provides the PostgreSQL implementation of the Repository interface for the Account service.
 * It encapsulates all database-specific logic, such as creating account records and fetching user data.
 *
 * @dependencies
 * - Go standard library packages: "context", "errors", "fmt"
 * - "github.com/jackc/pgx/v5": For checking specific database errors.
 * - "github.com/jackc/pgx/v5/pgxpool": The PostgreSQL driver and connection pool.
 * - "transfa/services/account/internal/domain": For core data models.
 */
package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"transfa/services/account/internal/domain"
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

// CreateAccount inserts a new account record into the public.accounts table.
func (r *PostgresRepository) CreateAccount(ctx context.Context, account *domain.Account) (*domain.Account, error) {
	query := `
        INSERT INTO public.accounts (user_id, anchor_account_id, account_purpose, balance, status)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(ctx, query,
		account.UserID,
		account.AnchorAccountID,
		account.AccountPurpose,
		account.Balance,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert account into database: %w", err)
	}

	return account, nil
}

// GetUserByID retrieves a user's ID and account type from the database.
// This is necessary to determine what kind of Anchor account to create.
func (r *PostgresRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	query := `SELECT id, account_type FROM public.users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, userID).Scan(&user.ID, &user.AccountType)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: with id %s", ErrUserNotFound, userID)
		}
		return nil, fmt.Errorf("failed to query user by id: %w", err)
	}

	return &user, nil
}
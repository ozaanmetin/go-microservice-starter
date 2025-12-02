package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

// Repository defines the interface for user data operations
type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

// repository implements the Repository interface using sqlx
type repository struct {
	db *sqlx.DB
}

// NewRepository creates a new user repository
func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

// Create inserts a new user into the database
func (r *repository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	err := r.db.QueryRowxContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Check for unique constraint violation (email already exists)
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			return ErrUserAlreadyExists
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetByID retrieves a user by their ID
func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by their email
func (r *repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Update updates an existing user
func (r *repository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, first_name = $3, last_name = $4, is_active = $5, updated_at = $6
		WHERE id = $7
	`

	user.UpdatedAt = time.Now().UTC()
	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete removes a user from the database
func (r *repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

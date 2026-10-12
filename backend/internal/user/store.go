package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/h-karlsson/uptime/internal/pgerr"
	"github.com/h-karlsson/uptime/internal/store/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	queries *db.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{queries: db.New(pool)}
}

func (store *PostgresStore) Create(ctx context.Context, email, passwordHash string) (User, error) {
	row, err := store.queries.CreateUser(ctx, db.CreateUserParams{Email: email, PasswordHash: passwordHash})
	if err != nil {
		if pgerr.IsUniqueViolation(err, "users_email_key") {
			return User{}, ErrEmailTaken
		}
		return User{}, fmt.Errorf("insert user: %w", err)
	}
	return User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, nil
}

// GetByEmail returns the user and their password hash, or ErrNotFound.
func (store *PostgresStore) GetByEmail(ctx context.Context, email string) (User, string, error) {
	row, err := store.queries.GetUserByEmail(ctx, email)
	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, "", ErrNotFound
	}
	if err != nil {
		return User{}, "", fmt.Errorf("get user by email: %w", err)
	}
	return User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, row.PasswordHash, nil
}

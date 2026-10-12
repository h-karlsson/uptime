package session

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/h-karlsson/uptime/internal/store/db"
	"github.com/h-karlsson/uptime/internal/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	queries *db.Queries
}

func NewPostgresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{queries: db.New(pool)}
}

func (store *PostgresStore) Create(ctx context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error {
	err := store.queries.CreateSession(ctx, db.CreateSessionParams{
		TokenHash: tokenHash,
		UserID:    userID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

// GetUser returns the user of an unexpired session, or ErrInvalidSession.
func (store *PostgresStore) GetUser(ctx context.Context, tokenHash []byte) (user.User, error) {
	row, err := store.queries.GetSessionUser(ctx, tokenHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, ErrInvalidSession
	}
	if err != nil {
		return user.User{}, fmt.Errorf("get session user: %w", err)
	}
	return user.User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt}, nil
}

func (store *PostgresStore) Delete(ctx context.Context, tokenHash []byte) error {
	if err := store.queries.DeleteSession(ctx, tokenHash); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

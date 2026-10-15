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

func (store *PostgresStore) Get(ctx context.Context, tokenHash []byte) (Session, error) {
	row, err := store.queries.GetSession(ctx, tokenHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return Session{}, ErrInvalidSession
	}
	if err != nil {
		return Session{}, fmt.Errorf("get session: %w", err)
	}
	return Session{
		User:      user.User{ID: row.ID, Email: row.Email, CreatedAt: row.CreatedAt},
		ExpiresAt: row.ExpiresAt,
	}, nil
}

func (store *PostgresStore) Extend(ctx context.Context, tokenHash []byte, expiresAt time.Time) error {
	err := store.queries.ExtendSession(ctx, db.ExtendSessionParams{TokenHash: tokenHash, ExpiresAt: expiresAt})
	if err != nil {
		return fmt.Errorf("extend session: %w", err)
	}
	return nil
}

func (store *PostgresStore) Delete(ctx context.Context, tokenHash []byte) error {
	if err := store.queries.DeleteSession(ctx, tokenHash); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}

func (store *PostgresStore) DeleteForUser(ctx context.Context, userID int64) error {
	if err := store.queries.DeleteUserSessions(ctx, userID); err != nil {
		return fmt.Errorf("delete user sessions: %w", err)
	}
	return nil
}

func (store *PostgresStore) DeleteExpired(ctx context.Context) (int64, error) {
	deleted, err := store.queries.DeleteExpiredSessions(ctx)
	if err != nil {
		return 0, fmt.Errorf("delete expired sessions: %w", err)
	}
	return deleted, nil
}

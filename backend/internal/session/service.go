package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"log/slog"
	"time"
)

type Store interface {
	Create(ctx context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error
	Get(ctx context.Context, tokenHash []byte) (Session, error)
	Extend(ctx context.Context, tokenHash []byte, expiresAt time.Time) error
	Delete(ctx context.Context, tokenHash []byte) error
	DeleteForUser(ctx context.Context, userID int64) error
	DeleteExpired(ctx context.Context) (int64, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

// Create starts a session for the user and returns the token to put in the cookie.
func (service *Service) Create(ctx context.Context, userID int64) (string, time.Time, error) {
	token := rand.Text() // 128 bits from crypto/rand
	expiresAt := time.Now().Add(Duration)
	if err := service.store.Create(ctx, hashToken(token), userID, expiresAt); err != nil {
		return "", time.Time{}, err
	}
	return token, expiresAt, nil
}

func (service *Service) Get(ctx context.Context, token string) (Session, error) {
	if token == "" {
		return Session{}, ErrInvalidSession
	}
	tokenHash := hashToken(token)
	current, err := service.store.Get(ctx, tokenHash)
	if err != nil {
		return Session{}, err
	}
	if time.Until(current.ExpiresAt) >= renewBelow {
		return current, nil
	}

	expiresAt := time.Now().Add(Duration)
	if err := service.store.Extend(ctx, tokenHash, expiresAt); err != nil {
		// The session is still valid, so the request goes ahead; the next one retries.
		slog.Warn("renew session", "user", current.User.ID, "error", err)
		return current, nil
	}
	current.ExpiresAt = expiresAt
	current.Renewed = true
	return current, nil
}

func (service *Service) Delete(ctx context.Context, token string) error {
	return service.store.Delete(ctx, hashToken(token))
}

func (service *Service) DeleteForUser(ctx context.Context, userID int64) error {
	return service.store.DeleteForUser(ctx, userID)
}

func hashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

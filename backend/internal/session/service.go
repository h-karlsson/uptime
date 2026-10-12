package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"time"

	"github.com/h-karlsson/uptime/internal/user"
)

type Store interface {
	Create(ctx context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error
	GetUser(ctx context.Context, tokenHash []byte) (user.User, error)
	Delete(ctx context.Context, tokenHash []byte) error
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

// User returns the logged-in user for a token, or ErrInvalidSession.
func (service *Service) User(ctx context.Context, token string) (user.User, error) {
	if token == "" {
		return user.User{}, ErrInvalidSession
	}
	return service.store.GetUser(ctx, hashToken(token))
}

func (service *Service) Delete(ctx context.Context, token string) error {
	return service.store.Delete(ctx, hashToken(token))
}

func hashToken(token string) []byte {
	hash := sha256.Sum256([]byte(token))
	return hash[:]
}

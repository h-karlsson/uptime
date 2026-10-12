package session

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/user"
)

type fakeStore struct {
	userIDs map[string]int64
}

func (store *fakeStore) Create(_ context.Context, tokenHash []byte, userID int64, _ time.Time) error {
	store.userIDs[string(tokenHash)] = userID
	return nil
}

func (store *fakeStore) GetUser(_ context.Context, tokenHash []byte) (user.User, error) {
	userID, exists := store.userIDs[string(tokenHash)]
	if !exists {
		return user.User{}, ErrInvalidSession
	}
	return user.User{ID: userID}, nil
}

func (store *fakeStore) Delete(_ context.Context, tokenHash []byte) error {
	delete(store.userIDs, string(tokenHash))
	return nil
}

func TestSessionLifecycle(t *testing.T) {
	store := &fakeStore{userIDs: map[string]int64{}}
	service := NewService(store)
	ctx := context.Background()

	token, expiresAt, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if time.Until(expiresAt) < Duration-time.Minute {
		t.Fatalf("expiresAt = %v, want about %v from now", expiresAt, Duration)
	}
	for storedHash := range store.userIDs {
		if bytes.Equal([]byte(storedHash), []byte(token)) {
			t.Fatal("raw token was stored; only its hash should be")
		}
	}

	loggedIn, err := service.User(ctx, token)
	if err != nil || loggedIn.ID != 42 {
		t.Fatalf("User = %+v, %v; want user 42", loggedIn, err)
	}
	if _, err := service.User(ctx, "not-a-real-token"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("unknown token: err = %v, want ErrInvalidSession", err)
	}
	if _, err := service.User(ctx, ""); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("empty token: err = %v, want ErrInvalidSession", err)
	}

	if err := service.Delete(ctx, token); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := service.User(ctx, token); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("after Delete: err = %v, want ErrInvalidSession", err)
	}
}

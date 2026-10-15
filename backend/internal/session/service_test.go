package session

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/user"
)

type fakeSession struct {
	userID    int64
	expiresAt time.Time
}

type fakeStore struct {
	sessions map[string]fakeSession
}

func newFakeStore() *fakeStore {
	return &fakeStore{sessions: map[string]fakeSession{}}
}

func (store *fakeStore) Create(_ context.Context, tokenHash []byte, userID int64, expiresAt time.Time) error {
	store.sessions[string(tokenHash)] = fakeSession{userID: userID, expiresAt: expiresAt}
	return nil
}

func (store *fakeStore) Get(_ context.Context, tokenHash []byte) (Session, error) {
	stored, exists := store.sessions[string(tokenHash)]
	if !exists || !stored.expiresAt.After(time.Now()) {
		return Session{}, ErrInvalidSession
	}
	return Session{User: user.User{ID: stored.userID}, ExpiresAt: stored.expiresAt}, nil
}

func (store *fakeStore) Extend(_ context.Context, tokenHash []byte, expiresAt time.Time) error {
	stored := store.sessions[string(tokenHash)]
	stored.expiresAt = expiresAt
	store.sessions[string(tokenHash)] = stored
	return nil
}

func (store *fakeStore) Delete(_ context.Context, tokenHash []byte) error {
	delete(store.sessions, string(tokenHash))
	return nil
}

func (store *fakeStore) DeleteForUser(_ context.Context, userID int64) error {
	for tokenHash, stored := range store.sessions {
		if stored.userID == userID {
			delete(store.sessions, tokenHash)
		}
	}
	return nil
}

func (store *fakeStore) DeleteExpired(context.Context) (int64, error) {
	return 0, nil
}

func TestSessionLifecycle(t *testing.T) {
	store := newFakeStore()
	service := NewService(store)
	ctx := context.Background()

	token, expiresAt, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if time.Until(expiresAt) < Duration-time.Minute {
		t.Fatalf("expiresAt = %v, want about %v from now", expiresAt, Duration)
	}
	for storedHash := range store.sessions {
		if bytes.Equal([]byte(storedHash), []byte(token)) {
			t.Fatal("raw token was stored; only its hash should be")
		}
	}

	current, err := service.Get(ctx, token)
	if err != nil || current.User.ID != 42 {
		t.Fatalf("Get = %+v, %v; want user 42", current, err)
	}
	if current.Renewed {
		t.Fatal("a fresh session should not be renewed")
	}
	if _, err := service.Get(ctx, "not-a-real-token"); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("unknown token: err = %v, want ErrInvalidSession", err)
	}
	if _, err := service.Get(ctx, ""); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("empty token: err = %v, want ErrInvalidSession", err)
	}

	if err := service.Delete(ctx, token); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := service.Get(ctx, token); !errors.Is(err, ErrInvalidSession) {
		t.Fatalf("after Delete: err = %v, want ErrInvalidSession", err)
	}
}

func TestGetRenewsSessionPastHalfLife(t *testing.T) {
	store := newFakeStore()
	service := NewService(store)
	ctx := context.Background()

	token, _, err := service.Create(ctx, 42)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Extend(ctx, hashToken(token), time.Now().Add(24*time.Hour)); err != nil {
		t.Fatal(err)
	}

	current, err := service.Get(ctx, token)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !current.Renewed || time.Until(current.ExpiresAt) < Duration-time.Minute {
		t.Fatalf("Get = %+v, want renewed to about %v from now", current, Duration)
	}
	if stored := store.sessions[string(hashToken(token))]; !stored.expiresAt.Equal(current.ExpiresAt) {
		t.Fatalf("stored expiry = %v, want %v", stored.expiresAt, current.ExpiresAt)
	}
}

func TestDeleteForUserEndsOnlyThatUsersSessions(t *testing.T) {
	service := NewService(newFakeStore())
	ctx := context.Background()

	laptop, _, _ := service.Create(ctx, 42)
	phone, _, _ := service.Create(ctx, 42)
	other, _, _ := service.Create(ctx, 7)

	if err := service.DeleteForUser(ctx, 42); err != nil {
		t.Fatalf("DeleteForUser: %v", err)
	}
	for _, token := range []string{laptop, phone} {
		if _, err := service.Get(ctx, token); !errors.Is(err, ErrInvalidSession) {
			t.Fatalf("user 42 session still valid: err = %v", err)
		}
	}
	if _, err := service.Get(ctx, other); err != nil {
		t.Fatalf("user 7 session was deleted: %v", err)
	}
}

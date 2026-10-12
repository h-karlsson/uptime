package user

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// fakeStore keeps users in memory and enforces unique emails like the database does.
type fakeStore struct {
	hashes map[string]string
}

func newFakeStore(existingEmails ...string) *fakeStore {
	store := &fakeStore{hashes: map[string]string{}}
	for _, email := range existingEmails {
		store.hashes[email] = "existing-hash"
	}
	return store
}

func (store *fakeStore) Create(_ context.Context, email, passwordHash string) (User, error) {
	if _, exists := store.hashes[email]; exists {
		return User{}, ErrEmailTaken
	}
	store.hashes[email] = passwordHash
	return User{ID: int64(len(store.hashes)), Email: email}, nil
}

func (store *fakeStore) GetByEmail(_ context.Context, email string) (User, string, error) {
	hash, exists := store.hashes[email]
	if !exists {
		return User{}, "", ErrNotFound
	}
	return User{Email: email}, hash, nil
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		want     error
	}{
		{"blank email", "  ", "longenough", ErrInvalidEmail},
		{"not an email", "henrik", "longenough", ErrInvalidEmail},
		{"display name", "Henrik <henrik@uptime.io>", "longenough", ErrInvalidEmail},
		{"short password", "henrik@uptime.io", "short", ErrPasswordTooShort},
		{"password over bcrypt limit", "henrik@uptime.io", strings.Repeat("a", MaxPasswordBytes+1), ErrPasswordTooLong},
		{"email taken", "demo@uptime.io", "longenough", ErrEmailTaken},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(newFakeStore("demo@uptime.io"))
			_, err := service.Register(context.Background(), test.email, test.password)
			if !errors.Is(err, test.want) {
				t.Fatalf("err = %v, want %v", err, test.want)
			}
		})
	}
}

func TestRegisterNormalizesEmailAndHashesPassword(t *testing.T) {
	store := newFakeStore()
	service := NewService(store)

	user, err := service.Register(context.Background(), "  Henrik@Uptime.IO ", "longenough")
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if user.Email != "henrik@uptime.io" {
		t.Fatalf("email = %q, want normalized address", user.Email)
	}

	hash := store.hashes["henrik@uptime.io"]
	if hash == "longenough" {
		t.Fatal("password was stored in plaintext")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte("longenough")); err != nil {
		t.Fatalf("stored hash does not match password: %v", err)
	}
}

func TestAuthenticate(t *testing.T) {
	service := NewService(newFakeStore())
	if _, err := service.Register(context.Background(), "henrik@uptime.io", "longenough"); err != nil {
		t.Fatalf("Register: %v", err)
	}

	tests := []struct {
		name     string
		email    string
		password string
		want     error
	}{
		{"correct", "henrik@uptime.io", "longenough", nil},
		{"email in other casing", " Henrik@Uptime.IO ", "longenough", nil},
		{"wrong password", "henrik@uptime.io", "wrongpassword", ErrInvalidCredentials},
		{"unknown email", "nobody@uptime.io", "longenough", ErrInvalidCredentials},
		{"invalid email", "henrik", "longenough", ErrInvalidCredentials},
		{"password over bcrypt limit", "henrik@uptime.io", strings.Repeat("a", MaxPasswordBytes+1), ErrInvalidCredentials},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := service.Authenticate(context.Background(), test.email, test.password)
			if !errors.Is(err, test.want) {
				t.Fatalf("err = %v, want %v", err, test.want)
			}
			if test.want == nil && user.Email != "henrik@uptime.io" {
				t.Fatalf("email = %q, want henrik@uptime.io", user.Email)
			}
		})
	}
}

func TestIsRegistered(t *testing.T) {
	service := NewService(newFakeStore("demo@uptime.io"))

	tests := []struct {
		name    string
		email   string
		want    bool
		wantErr error
	}{
		{"registered", "demo@uptime.io", true, nil},
		{"registered in other casing", " Demo@Uptime.IO ", true, nil},
		{"not registered", "henrik@uptime.io", false, nil},
		{"invalid email", "henrik", false, ErrInvalidEmail},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			registered, err := service.IsRegistered(context.Background(), test.email)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("err = %v, want %v", err, test.wantErr)
			}
			if registered != test.want {
				t.Fatalf("registered = %v, want %v", registered, test.want)
			}
		})
	}
}

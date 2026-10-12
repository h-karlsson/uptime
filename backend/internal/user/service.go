package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Store interface {
	Create(ctx context.Context, email, passwordHash string) (User, error)
	GetByEmail(ctx context.Context, email string) (User, string, error)
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) Register(ctx context.Context, email, password string) (User, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return User{}, err
	}
	if err := validatePassword(password); err != nil {
		return User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}
	return service.store.Create(ctx, email, string(hash))
}

// Authenticate returns ErrInvalidCredentials for an unknown email or a wrong password,
// so login responses don't reveal which of the two was wrong.
func (service *Service) Authenticate(ctx context.Context, email, password string) (User, error) {
	email, err := normalizeEmail(email)
	if err != nil || len(password) > MaxPasswordBytes {
		return User{}, ErrInvalidCredentials
	}

	user, passwordHash, err := service.store.GetByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return User{}, ErrInvalidCredentials
	}
	if err != nil {
		return User{}, fmt.Errorf("compare password: %w", err)
	}
	return user, nil
}

// IsRegistered backs the email-first login step: the client asks whether to show
// the password field or the registration form.
func (service *Service) IsRegistered(ctx context.Context, email string) (bool, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return false, err
	}

	_, _, err = service.store.GetByEmail(ctx, email)
	if errors.Is(err, ErrNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", ErrInvalidEmail
	}
	return email, nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}
	if len(password) > MaxPasswordBytes {
		return ErrPasswordTooLong
	}
	return nil
}

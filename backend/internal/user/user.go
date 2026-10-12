package user

import (
	"errors"
	"fmt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

const (
	MinPasswordLength = 8
	MaxPasswordBytes  = 72
)

type ValidationError struct {
	Message string
}

func (err *ValidationError) Error() string {
	return err.Message
}

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email is already registered")
	ErrInvalidEmail       = &ValidationError{Message: "email must be a valid address"}
	ErrPasswordTooShort   = &ValidationError{Message: fmt.Sprintf("password must be at least %d characters", MinPasswordLength)}
	ErrPasswordTooLong    = &ValidationError{Message: fmt.Sprintf("password must be at most %d bytes", MaxPasswordBytes)}
)

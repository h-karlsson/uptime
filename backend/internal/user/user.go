package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/h-karlsson/uptime/internal/validation"
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

var (
	ErrNotFound           = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email is already registered")
	ErrInvalidEmail       = &validation.Error{Message: "email must be a valid address"}
	ErrPasswordTooShort   = &validation.Error{Message: fmt.Sprintf("password must be at least %d characters", MinPasswordLength)}
	ErrPasswordTooLong    = &validation.Error{Message: fmt.Sprintf("password must be at most %d bytes", MaxPasswordBytes)}
)

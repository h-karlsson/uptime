package session

import (
	"errors"
	"time"

	"github.com/h-karlsson/uptime/internal/user"
)

const (
	CookieName = "session"
	Duration   = 30 * 24 * time.Hour
	renewBelow = Duration / 2
)

var ErrInvalidSession = errors.New("session is invalid or expired")

type Session struct {
	User      user.User
	ExpiresAt time.Time
	Renewed   bool
}

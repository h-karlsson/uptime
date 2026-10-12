package session

import (
	"errors"
	"time"
)

const (
	CookieName = "session"
	Duration   = 30 * 24 * time.Hour
)

var ErrInvalidSession = errors.New("session is invalid or expired")

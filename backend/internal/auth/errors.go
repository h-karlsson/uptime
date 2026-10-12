package auth

import (
	"errors"
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
)

var errNotLoggedIn = &customHttp.ClientError{Message: "not logged in", Status: http.StatusUnauthorized}

// toClientError turns domain errors the client should see into HTTP errors; anything else stays a 500.
func toClientError(err error) error {
	switch {
	case errors.Is(err, user.ErrEmailTaken):
		return &customHttp.ClientError{Message: err.Error(), Status: http.StatusConflict}
	case errors.Is(err, user.ErrInvalidCredentials):
		return &customHttp.ClientError{Message: err.Error(), Status: http.StatusUnauthorized}
	case errors.Is(err, session.ErrInvalidSession):
		return errNotLoggedIn
	}
	if validationError, ok := errors.AsType[*user.ValidationError](err); ok {
		return &customHttp.ClientError{Message: validationError.Message}
	}
	return err
}

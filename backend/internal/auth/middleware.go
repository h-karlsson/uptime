package auth

import (
	"context"
	"errors"
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
)

type contextKey struct{}

// RequireUser responds 401 unless the request has a valid session cookie. Handlers behind
// it get the logged-in user with CurrentUser.
func RequireUser(sessions *session.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			loggedIn, err := sessions.User(request.Context(), sessionToken(request))
			if err != nil {
				customHttp.WriteError(response, toClientError(err))
				return
			}
			next.ServeHTTP(response, request.WithContext(WithUser(request.Context(), loggedIn)))
		})
	}
}

func WithUser(ctx context.Context, loggedIn user.User) context.Context {
	return context.WithValue(ctx, contextKey{}, loggedIn)
}

// CurrentUser returns the logged-in user. It only fails when the route is missing RequireUser.
func CurrentUser(ctx context.Context) (user.User, error) {
	loggedIn, ok := ctx.Value(contextKey{}).(user.User)
	if !ok {
		return user.User{}, errors.New("no user in context; route is missing RequireUser")
	}
	return loggedIn, nil
}

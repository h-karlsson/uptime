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

func RequireUser(sessions *session.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			token := sessionToken(request)
			current, err := sessions.Get(request.Context(), token)
			if err != nil {
				customHttp.WriteError(response, toClientError(err))
				return
			}
			if current.Renewed {
				setSessionCookie(response, token, current.ExpiresAt)
			}
			next.ServeHTTP(response, request.WithContext(WithUser(request.Context(), current.User)))
		})
	}
}

func WithUser(ctx context.Context, loggedIn user.User) context.Context {
	return context.WithValue(ctx, contextKey{}, loggedIn)
}

func CurrentUser(ctx context.Context) (user.User, error) {
	loggedIn, ok := ctx.Value(contextKey{}).(user.User)
	if !ok {
		return user.User{}, errors.New("no user in context; route is missing RequireUser")
	}
	return loggedIn, nil
}

func CurrentUserID(ctx context.Context) (int64, error) {
	loggedIn, err := CurrentUser(ctx)
	return loggedIn.ID, err
}

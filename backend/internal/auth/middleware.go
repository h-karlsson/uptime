package auth

import (
	"context"
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
)

type contextKey struct{}

// RequireUser responds 401 unless the request has a valid session cookie. Handlers behind
// it get the logged-in user with UserFromContext.
func RequireUser(sessions *session.Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			loggedIn, err := sessions.User(request.Context(), sessionToken(request))
			if err != nil {
				customHttp.WriteError(response, toClientError(err))
				return
			}
			ctx := context.WithValue(request.Context(), contextKey{}, loggedIn)
			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (user.User, bool) {
	loggedIn, ok := ctx.Value(contextKey{}).(user.User)
	return loggedIn, ok
}

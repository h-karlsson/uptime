package auth

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
)

// Logout deletes the session and clears the cookie. It succeeds even without a session,
// so the client can always call it.
func Logout(sessions *session.Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		clearSessionCookie(response)
		if token := sessionToken(request); token != "" {
			if err := sessions.Delete(request.Context(), token); err != nil {
				customHttp.WriteError(response, err)
				return
			}
		}
		response.WriteHeader(http.StatusNoContent)
	}
}

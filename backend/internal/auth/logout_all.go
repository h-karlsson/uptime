package auth

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
)

func LogoutAll(sessions *session.Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		userID, err := CurrentUserID(request.Context())
		if err != nil {
			customHttp.WriteError(response, err)
			return
		}
		if err := sessions.DeleteForUser(request.Context(), userID); err != nil {
			customHttp.WriteError(response, err)
			return
		}
		clearSessionCookie(response)
		response.WriteHeader(http.StatusNoContent)
	}
}

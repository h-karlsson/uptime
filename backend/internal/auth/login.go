package auth

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(users *user.Service, sessions *session.Service) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var requestBody LoginRequest
		if err := customHttp.ReadJSON(request, &requestBody); err != nil {
			customHttp.WriteError(response, err)
			return
		}

		loggedIn, err := users.Authenticate(request.Context(), requestBody.Email, requestBody.Password)
		if err != nil {
			customHttp.WriteError(response, toClientError(err))
			return
		}
		if err := startSession(request.Context(), response, sessions, loggedIn.ID); err != nil {
			customHttp.WriteError(response, err)
			return
		}
		customHttp.FormatResponseToJson(response, http.StatusOK, loggedIn)
	}
}

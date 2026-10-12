package auth

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/user"
)

type CheckEmailRequest struct {
	Email string `json:"email"`
}

type CheckEmailResponse struct {
	Registered bool `json:"registered"`
}

func CheckEmail(users *user.Service) func(*http.Request) (CheckEmailResponse, error) {
	return func(request *http.Request) (CheckEmailResponse, error) {
		var requestBody CheckEmailRequest
		if err := customHttp.ReadJSON(request, &requestBody); err != nil {
			return CheckEmailResponse{}, err
		}

		registered, err := users.IsRegistered(request.Context(), requestBody.Email)
		if err != nil {
			return CheckEmailResponse{}, toClientError(err)
		}
		return CheckEmailResponse{Registered: registered}, nil
	}
}

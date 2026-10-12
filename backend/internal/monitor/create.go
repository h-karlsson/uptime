package monitor

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func Create(monitors *Service) func(*http.Request) (Monitor, error) {
	return func(request *http.Request) (Monitor, error) {
		userID, err := auth.CurrentUserID(request.Context())
		if err != nil {
			return Monitor{}, err
		}
		var input Input
		if err := customHttp.ReadJSON(request, &input); err != nil {
			return Monitor{}, err
		}

		created, err := monitors.Create(request.Context(), userID, input)
		if err != nil {
			return Monitor{}, toClientError(err)
		}
		return created, nil
	}
}

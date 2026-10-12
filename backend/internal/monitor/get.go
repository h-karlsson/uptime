package monitor

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
)

func Get(monitors *Service) func(*http.Request) (Monitor, error) {
	return func(request *http.Request) (Monitor, error) {
		userID, err := auth.CurrentUserID(request.Context())
		if err != nil {
			return Monitor{}, err
		}
		id, err := monitorID(request)
		if err != nil {
			return Monitor{}, toClientError(err)
		}

		found, err := monitors.Get(request.Context(), userID, id)
		if err != nil {
			return Monitor{}, toClientError(err)
		}
		return found, nil
	}
}

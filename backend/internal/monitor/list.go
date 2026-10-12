package monitor

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
)

func List(monitors *Service) func(*http.Request) ([]Monitor, error) {
	return func(request *http.Request) ([]Monitor, error) {
		userID, err := auth.CurrentUserID(request.Context())
		if err != nil {
			return nil, err
		}
		return monitors.List(request.Context(), userID)
	}
}

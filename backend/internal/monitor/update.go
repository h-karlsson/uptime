package monitor

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func Update(monitors *Service) func(*http.Request) (Monitor, error) {
	return func(request *http.Request) (Monitor, error) {
		userID, err := auth.CurrentUserID(request.Context())
		if err != nil {
			return Monitor{}, err
		}
		id, err := monitorID(request)
		if err != nil {
			return Monitor{}, toClientError(err)
		}
		var input Input
		if err := customHttp.ReadJSON(request, &input); err != nil {
			return Monitor{}, err
		}

		updated, err := monitors.Update(request.Context(), userID, id, input)
		if err != nil {
			return Monitor{}, toClientError(err)
		}
		return updated, nil
	}
}

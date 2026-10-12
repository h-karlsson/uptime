package monitor

import "net/http"

func List(monitors *Service) func(*http.Request) ([]Monitor, error) {
	return func(request *http.Request) ([]Monitor, error) {
		userID, err := currentUserID(request)
		if err != nil {
			return nil, err
		}
		return monitors.List(request.Context(), userID)
	}
}

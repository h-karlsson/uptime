package monitor

import "net/http"

func Get(monitors *Service) func(*http.Request) (Monitor, error) {
	return func(request *http.Request) (Monitor, error) {
		userID, err := currentUserID(request)
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

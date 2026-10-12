package monitor

import "net/http"

func Delete(monitors *Service) func(*http.Request) error {
	return func(request *http.Request) error {
		userID, err := currentUserID(request)
		if err != nil {
			return err
		}
		id, err := monitorID(request)
		if err != nil {
			return toClientError(err)
		}
		return toClientError(monitors.Delete(request.Context(), userID, id))
	}
}

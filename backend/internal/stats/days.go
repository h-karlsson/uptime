package stats

import "net/http"

func ListDays(stats *Service) func(*http.Request) ([]Day, error) {
	return func(request *http.Request) ([]Day, error) {
		userID, monitorID, err := requestIDs(request)
		if err != nil {
			return nil, err
		}
		days, err := stats.Days(request.Context(), userID, monitorID)
		return days, toClientError(err)
	}
}

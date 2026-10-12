package stats

import "net/http"

func GetSummary(stats *Service) func(*http.Request) (Summary, error) {
	return func(request *http.Request) (Summary, error) {
		userID, monitorID, err := requestIDs(request)
		if err != nil {
			return Summary{}, err
		}
		summary, err := stats.Summary(request.Context(), userID, monitorID, request.URL.Query().Get("range"))
		return summary, toClientError(err)
	}
}

package stats

import "net/http"

func ListIncidents(stats *Service) func(*http.Request) ([]Incident, error) {
	return func(request *http.Request) ([]Incident, error) {
		userID, monitorID, err := requestIDs(request)
		if err != nil {
			return nil, err
		}
		incidents, err := stats.Incidents(request.Context(), userID, monitorID)
		return incidents, toClientError(err)
	}
}

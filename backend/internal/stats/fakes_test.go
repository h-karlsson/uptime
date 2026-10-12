package stats

import (
	"context"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/h-karlsson/uptime/internal/auth"
	"github.com/h-karlsson/uptime/internal/user"
)

// fakeStore owns monitor 1 for user 1 and records the time windows it was asked for.
type fakeStore struct {
	since  time.Time
	bucket time.Duration
}

func (store *fakeStore) BelongsToUser(_ context.Context, userID, monitorID int64) (bool, error) {
	return userID == 1 && monitorID == 1, nil
}

func (store *fakeStore) Summary(_ context.Context, _ int64, since time.Time) (Summary, error) {
	store.since = since
	return Summary{Checks: 1}, nil
}

func (store *fakeStore) Latency(_ context.Context, _ int64, _ time.Time, bucket time.Duration) ([]LatencyPoint, error) {
	store.bucket = bucket
	return []LatencyPoint{}, nil
}

func (store *fakeStore) Days(_ context.Context, _ int64, since time.Time) ([]Day, error) {
	store.since = since
	return []Day{}, nil
}

func (store *fakeStore) Incidents(context.Context, int64) ([]Incident, error) {
	return []Incident{}, nil
}

func newTestRouter(store *fakeStore, userID int64) *http.ServeMux {
	loggedInAs := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			next.ServeHTTP(response, request.WithContext(auth.WithUser(request.Context(), user.User{ID: userID})))
		})
	}
	router := http.NewServeMux()
	RegisterRoutes(router, NewService(store), loggedInAs)
	return router
}

func send(router http.Handler, path string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response
}

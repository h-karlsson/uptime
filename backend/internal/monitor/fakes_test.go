package monitor

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/h-karlsson/uptime/internal/auth"
	"github.com/h-karlsson/uptime/internal/user"
)

type ownedMonitor struct {
	userID  int64
	monitor Monitor
}

// fakeStore mimics PostgresStore, including scoping every lookup to the owner and the
// unique URL per user.
type fakeStore struct {
	monitors []ownedMonitor
	nextID   int64
}

func (store *fakeStore) List(_ context.Context, userID int64) ([]Monitor, error) {
	found := []Monitor{}
	for _, owned := range store.monitors {
		if owned.userID == userID {
			found = append(found, owned.monitor)
		}
	}
	return found, nil
}

func (store *fakeStore) Get(_ context.Context, userID, id int64) (Monitor, error) {
	index := store.find(userID, id)
	if index < 0 {
		return Monitor{}, ErrNotFound
	}
	return store.monitors[index].monitor, nil
}

func (store *fakeStore) Count(ctx context.Context, userID int64) (int, error) {
	monitors, err := store.List(ctx, userID)
	return len(monitors), err
}

func (store *fakeStore) Create(_ context.Context, userID int64, input Input) (Monitor, error) {
	if store.hasURL(userID, 0, input.URL) {
		return Monitor{}, ErrDuplicateURL
	}
	store.nextID++
	created := fromInput(store.nextID, input)
	store.monitors = append(store.monitors, ownedMonitor{userID: userID, monitor: created})
	return created, nil
}

func (store *fakeStore) Update(_ context.Context, userID, id int64, input Input) (Monitor, error) {
	index := store.find(userID, id)
	if index < 0 {
		return Monitor{}, ErrNotFound
	}
	if store.hasURL(userID, id, input.URL) {
		return Monitor{}, ErrDuplicateURL
	}
	store.monitors[index].monitor = fromInput(id, input)
	return store.monitors[index].monitor, nil
}

func (store *fakeStore) Delete(_ context.Context, userID, id int64) error {
	index := store.find(userID, id)
	if index < 0 {
		return ErrNotFound
	}
	store.monitors = append(store.monitors[:index], store.monitors[index+1:]...)
	return nil
}

func (store *fakeStore) find(userID, id int64) int {
	for index, owned := range store.monitors {
		if owned.userID == userID && owned.monitor.ID == id {
			return index
		}
	}
	return -1
}

func (store *fakeStore) hasURL(userID, exceptID int64, url string) bool {
	for _, owned := range store.monitors {
		if owned.userID == userID && owned.monitor.ID != exceptID && owned.monitor.URL == url {
			return true
		}
	}
	return false
}

func fromInput(id int64, input Input) Monitor {
	return Monitor{
		ID:              id,
		Name:            input.Name,
		URL:             input.URL,
		Method:          input.Method,
		IntervalSeconds: input.IntervalSeconds,
		TimeoutMs:       input.TimeoutMs,
		ExpectedStatus:  input.ExpectedStatus,
		Enabled:         input.Enabled,
		Status:          "unknown",
	}
}

func validInput() Input {
	return Input{
		Name:            "Example",
		URL:             "https://example.com",
		Method:          "GET",
		IntervalSeconds: 60,
		TimeoutMs:       10000,
		ExpectedStatus:  200,
		Enabled:         true,
	}
}

// newTestRouter registers the routes with a stand-in for RequireUser that logs in as userID.
func newTestRouter(store *fakeStore, userID int64) *http.ServeMux {
	loggedInAs := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			ctx := auth.WithUser(request.Context(), user.User{ID: userID})
			next.ServeHTTP(response, request.WithContext(ctx))
		})
	}
	router := http.NewServeMux()
	RegisterRoutes(router, NewService(store), loggedInAs)
	return router
}

func send(router http.Handler, method, path, body string) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(method, path, strings.NewReader(body)))
	return response
}

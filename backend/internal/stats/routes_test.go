package stats

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		name   string
		userID int64
		path   string
		status int
	}{
		{"stats", 1, "/api/monitors/1/stats?range=7d", http.StatusOK},
		{"days", 1, "/api/monitors/1/days", http.StatusOK},
		{"incidents", 1, "/api/monitors/1/incidents", http.StatusOK},
		{"unknown range", 1, "/api/monitors/1/stats?range=1y", http.StatusBadRequest},
		{"non-numeric id", 1, "/api/monitors/abc/stats", http.StatusNotFound},
		{"another user's stats", 2, "/api/monitors/1/stats", http.StatusNotFound},
		{"another user's days", 2, "/api/monitors/1/days", http.StatusNotFound},
		{"another user's incidents", 2, "/api/monitors/1/incidents", http.StatusNotFound},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := send(newTestRouter(&fakeStore{}, test.userID), test.path)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}

func TestEmptyListsAreJSONArrays(t *testing.T) {
	router := newTestRouter(&fakeStore{}, 1)
	var summary map[string]any
	if err := json.Unmarshal(send(router, "/api/monitors/1/stats").Body.Bytes(), &summary); err != nil {
		t.Fatal(err)
	}
	if _, ok := summary["latency"].([]any); !ok {
		t.Fatalf("latency = %v, want an empty array", summary["latency"])
	}
	for _, path := range []string{"/api/monitors/1/days", "/api/monitors/1/incidents"} {
		if body := send(router, path).Body.String(); body != "[]\n" {
			t.Fatalf("%s: body = %q, want []", path, body)
		}
	}
}

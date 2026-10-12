package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakePinger struct {
	err error
}

func (pinger fakePinger) Ping(context.Context) error {
	return pinger.err
}

func TestGetHealth(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{"database up", nil, http.StatusOK},
		{"database down", errors.New("connection refused"), http.StatusServiceUnavailable},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := http.NewServeMux()
			RegisterRoutes(router, fakePinger{err: test.err})
			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/api/health", nil)
			router.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
		})
	}
}

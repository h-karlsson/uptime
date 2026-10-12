package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTarget(url string) target {
	return target{monitorID: 1, url: url, method: http.MethodGet, timeout: time.Second, expectedStatus: http.StatusOK}
}

func TestProbe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		switch request.URL.Path {
		case "/ok":
			response.WriteHeader(http.StatusOK)
		case "/error":
			response.WriteHeader(http.StatusServiceUnavailable)
		case "/redirect":
			http.Redirect(response, request, "/ok", http.StatusFound)
		case "/loop":
			http.Redirect(response, request, "/loop", http.StatusFound)
		case "/slow":
			time.Sleep(200 * time.Millisecond)
		case "/head":
			if request.Method != http.MethodHead {
				response.WriteHeader(http.StatusMethodNotAllowed)
			}
		}
	}))
	defer server.Close()
	client := server.Client()
	client.CheckRedirect = newClient().CheckRedirect

	tests := []struct {
		name      string
		target    func(target) target
		ok        bool
		status    int32
		errSubstr string
	}{
		{"expected status", func(t target) target { t.url += "/ok"; return t }, true, 200, ""},
		{"unexpected status", func(t target) target { t.url += "/error"; return t }, false, 503, "expected status 200, got 503"},
		{"other expected status", func(t target) target { t.url += "/error"; t.expectedStatus = 503; return t }, true, 503, ""},
		{"follows redirects", func(t target) target { t.url += "/redirect"; return t }, true, 200, ""},
		{"too many redirects", func(t target) target { t.url += "/loop"; return t }, false, 0, "stopped after 5 redirects"},
		{"timeout", func(t target) target { t.url += "/slow"; t.timeout = 50 * time.Millisecond; return t }, false, 0, "timed out after 50ms"},
		{"head method", func(t target) target { t.url += "/head"; t.method = http.MethodHead; return t }, true, 200, ""},
		{"connection refused", func(t target) target { t.url = "http://127.0.0.1:1"; return t }, false, 0, "refused"},
		{"unknown host", func(t target) target { t.url = "http://no-such-host.invalid"; return t }, false, 0, "could not resolve host no-such-host.invalid"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := probe(context.Background(), client, test.target(newTarget(server.URL)))
			if result.ok != test.ok {
				t.Fatalf("ok = %v, want %v (err: %v)", result.ok, test.ok, deref(result.err))
			}
			if test.status != 0 && (result.statusCode == nil || *result.statusCode != test.status) {
				t.Fatalf("status = %v, want %d", deref(result.statusCode), test.status)
			}
			if test.ok && (result.err != nil || result.latencyMs == nil) {
				t.Fatalf("successful check: err = %v, latency = %v", deref(result.err), deref(result.latencyMs))
			}
			if test.errSubstr != "" && (result.err == nil || !strings.Contains(*result.err, test.errSubstr)) {
				t.Fatalf("err = %q, want it to contain %q", deref(result.err), test.errSubstr)
			}
		})
	}
}

func deref[T any](value *T) any {
	if value == nil {
		return nil
	}
	return *value
}

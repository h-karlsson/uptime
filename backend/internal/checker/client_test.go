package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/netip"
	"strings"
	"testing"
)

func TestIsPublic(t *testing.T) {
	tests := []struct {
		address string
		public  bool
	}{
		{"93.184.215.14", true},
		{"2606:2800:21f:cb07:6820:80da:af6b:8b2c", true},
		{"127.0.0.1", false},
		{"::1", false},
		{"10.0.0.5", false},
		{"172.16.3.4", false},
		{"192.168.1.1", false},
		{"169.254.169.254", false},
		{"100.64.0.1", false},
		{"0.0.0.0", false},
		{"224.0.0.1", false},
		{"fc00::1", false},
		{"fe80::1", false},
		{"::ffff:127.0.0.1", false},
	}
	for _, test := range tests {
		t.Run(test.address, func(t *testing.T) {
			if got := isPublic(netip.MustParseAddr(test.address)); got != test.public {
				t.Fatalf("isPublic(%s) = %v, want %v", test.address, got, test.public)
			}
		})
	}
}

func TestClientBlocksPrivateAddresses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Error("request reached a loopback server")
	}))
	defer server.Close()

	for _, url := range []string{server.URL, strings.Replace(server.URL, "127.0.0.1", "localhost", 1)} {
		result := probe(context.Background(), newClient(), newTarget(url))
		if result.ok || result.err == nil || !strings.Contains(*result.err, "blocked") {
			t.Fatalf("%s: err = %v, want the connection to be blocked", url, deref(result.err))
		}
	}
}

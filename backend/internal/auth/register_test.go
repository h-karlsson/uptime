package auth

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"empty", "", 400},
		{"whitespace", "  ", 400},
		{"malformed", `{`, 400},
		{"unknown field", `{"email":"henrik@uptime.io","password":"longenough","unexpected":true}`, 400},
		{"wrong type", `{"email":"henrik@uptime.io","password":12345678}`, 400},
		{"missing email", `{"password":"longenough"}`, 400},
		{"missing password", `{"email":"henrik@uptime.io"}`, 400},
		{"invalid email", `{"email":"henrik","password":"longenough"}`, 400},
		{"short password", `{"email":"henrik@uptime.io","password":"short"}`, 400},
		{"empty object", `{}`, 400},
		{"null", `null`, 400},
		{"array", `[]`, 400},
		{"multiple values", `{"email":"henrik@uptime.io","password":"longenough"} {}`, 400},
		{"email taken", `{"email":"demo@uptime.io","password":"longenough"}`, 409},
		{"valid", `{"email":"henrik@uptime.io","password":"longenough"}`, 201},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := send(newTestRouter(), http.MethodPost, "/api/auth/register", test.body)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}

			var body map[string]any
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatalf("expected JSON body, got %s", response.Body.String())
			}
			if test.status != 201 {
				if message, _ := body["error"].(string); message == "" || message == "Internal server error" {
					t.Fatalf("expected explanatory JSON error, got %s", response.Body.String())
				}
				if sessionCookie(response) != nil {
					t.Fatal("failed registration must not start a session")
				}
				return
			}
			if body["email"] != "henrik@uptime.io" {
				t.Fatalf("expected created user in response, got %s", response.Body.String())
			}
			if strings.Contains(strings.ToLower(response.Body.String()), "password") {
				t.Fatalf("response must not contain the password or its hash: %s", response.Body.String())
			}
			if sessionCookie(response) == nil {
				t.Fatal("registration should log the user in with a session cookie")
			}
		})
	}
}

package auth

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestCheckEmail(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		status     int
		registered bool
	}{
		{"registered", `{"email":"demo@uptime.io"}`, 200, true},
		{"registered in other casing", `{"email":" Demo@Uptime.IO "}`, 200, true},
		{"not registered", `{"email":"henrik@uptime.io"}`, 200, false},
		{"invalid email", `{"email":"henrik"}`, 400, false},
		{"malformed", `{`, 400, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := send(newTestRouter(), http.MethodPost, "/api/auth/check-email", test.body)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
			if test.status != 200 {
				return
			}
			var body CheckEmailResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || body.Registered != test.registered {
				t.Fatalf("body = %s, want registered = %v", response.Body.String(), test.registered)
			}
		})
	}
}

func TestLoginFailures(t *testing.T) {
	tests := []struct {
		name   string
		body   string
		status int
	}{
		{"wrong password", `{"email":"demo@uptime.io","password":"wrongpassword"}`, 401},
		{"unknown email", `{"email":"nobody@uptime.io","password":"demo1234"}`, 401},
		{"malformed", `{`, 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := send(newTestRouter(), http.MethodPost, "/api/auth/login", test.body)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.status, response.Body.String())
			}
			if sessionCookie(response) != nil {
				t.Fatal("failed login must not start a session")
			}
		})
	}
}

func TestLoginMeLogout(t *testing.T) {
	router := newTestRouter()

	if response := send(router, http.MethodGet, "/api/auth/me", ""); response.Code != 401 {
		t.Fatalf("me without cookie: status = %d, want 401", response.Code)
	}

	login := send(router, http.MethodPost, "/api/auth/login", `{"email":"Demo@Uptime.io","password":"demo1234"}`)
	if login.Code != 200 {
		t.Fatalf("login: status = %d, want 200; body = %s", login.Code, login.Body.String())
	}
	cookie := sessionCookie(login)
	if cookie == nil {
		t.Fatal("login should set a session cookie")
	}
	if !cookie.HttpOnly || !cookie.Secure || cookie.SameSite != http.SameSiteLaxMode {
		t.Fatalf("cookie must be HttpOnly, Secure and SameSite=Lax, got %+v", cookie)
	}

	me := send(router, http.MethodGet, "/api/auth/me", "", cookie)
	if me.Code != 200 {
		t.Fatalf("me: status = %d, want 200; body = %s", me.Code, me.Body.String())
	}
	var body map[string]any
	if err := json.Unmarshal(me.Body.Bytes(), &body); err != nil || body["email"] != "demo@uptime.io" {
		t.Fatalf("me: body = %s, want demo@uptime.io", me.Body.String())
	}

	logout := send(router, http.MethodPost, "/api/auth/logout", "", cookie)
	if logout.Code != 204 {
		t.Fatalf("logout: status = %d, want 204", logout.Code)
	}
	if cleared := sessionCookie(logout); cleared == nil || cleared.MaxAge >= 0 {
		t.Fatalf("logout should expire the cookie, got %+v", cleared)
	}

	if response := send(router, http.MethodGet, "/api/auth/me", "", cookie); response.Code != 401 {
		t.Fatalf("me after logout: status = %d, want 401", response.Code)
	}
}

func TestRegisterThenMe(t *testing.T) {
	router := newTestRouter()

	register := send(router, http.MethodPost, "/api/auth/register", `{"email":"henrik@uptime.io","password":"longenough"}`)
	cookie := sessionCookie(register)
	if register.Code != 201 || cookie == nil {
		t.Fatalf("register: status = %d, cookie = %v; want 201 with cookie", register.Code, cookie)
	}

	me := send(router, http.MethodGet, "/api/auth/me", "", cookie)
	var body map[string]any
	if err := json.Unmarshal(me.Body.Bytes(), &body); err != nil || me.Code != 200 || body["email"] != "henrik@uptime.io" {
		t.Fatalf("me: status = %d, body = %s; want 200 henrik@uptime.io", me.Code, me.Body.String())
	}
}

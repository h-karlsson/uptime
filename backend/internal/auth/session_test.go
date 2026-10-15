package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/session"
)

const demoLogin = `{"email":"demo@uptime.io","password":"demo1234"}`

func TestRequireUserRenewsAgingSession(t *testing.T) {
	router, sessions := newTestRouterWithSessions()
	cookie := sessionCookie(send(router, http.MethodPost, "/api/auth/login", demoLogin))

	if renewed := sessionCookie(send(router, http.MethodGet, "/api/auth/me", "", cookie)); renewed != nil {
		t.Fatalf("fresh session should not resend the cookie, got %+v", renewed)
	}

	sessions.ageSessions(24 * time.Hour)
	me := send(router, http.MethodGet, "/api/auth/me", "", cookie)
	renewed := sessionCookie(me)
	if me.Code != 200 || renewed == nil {
		t.Fatalf("me: status = %d, cookie = %v; want 200 with a renewed cookie", me.Code, renewed)
	}
	if renewed.Value != cookie.Value || time.Until(renewed.Expires) < session.Duration-time.Minute {
		t.Fatalf("renewed cookie = %+v, want same token expiring about %v from now", renewed, session.Duration)
	}
}

func TestLogoutAll(t *testing.T) {
	router := newTestRouter()
	laptop := sessionCookie(send(router, http.MethodPost, "/api/auth/login", demoLogin))
	phone := sessionCookie(send(router, http.MethodPost, "/api/auth/login", demoLogin))

	if response := send(router, http.MethodPost, "/api/auth/logout-all", ""); response.Code != 401 {
		t.Fatalf("logout-all without cookie: status = %d, want 401", response.Code)
	}

	logoutAll := send(router, http.MethodPost, "/api/auth/logout-all", "", laptop)
	if logoutAll.Code != 204 {
		t.Fatalf("logout-all: status = %d, want 204; body = %s", logoutAll.Code, logoutAll.Body.String())
	}
	if cleared := sessionCookie(logoutAll); cleared == nil || cleared.MaxAge >= 0 {
		t.Fatalf("logout-all should expire the cookie, got %+v", cleared)
	}

	for name, cookie := range map[string]*http.Cookie{"laptop": laptop, "phone": phone} {
		if response := send(router, http.MethodGet, "/api/auth/me", "", cookie); response.Code != 401 {
			t.Fatalf("me with %s cookie after logout-all: status = %d, want 401", name, response.Code)
		}
	}
}

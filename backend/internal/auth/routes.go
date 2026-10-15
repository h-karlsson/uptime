package auth

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
)

func RegisterRoutes(router *http.ServeMux, users *user.Service, sessions *session.Service) {
	requireUser := RequireUser(sessions)

	router.HandleFunc("POST /api/auth/check-email", customHttp.HandleJSON(CheckEmail(users)))
	router.HandleFunc("POST /api/auth/register", Register(users, sessions))
	router.HandleFunc("POST /api/auth/login", Login(users, sessions))
	router.HandleFunc("POST /api/auth/logout", Logout(sessions))
	router.Handle("POST /api/auth/logout-all", requireUser(LogoutAll(sessions)))
	router.Handle("GET /api/auth/me", requireUser(customHttp.HandleJSON(Me)))
}

package auth

import (
	"context"
	"net/http"

	"github.com/h-karlsson/uptime/internal/session"
)

func startSession(ctx context.Context, response http.ResponseWriter, sessions *session.Service, userID int64) error {
	token, expiresAt, err := sessions.Create(ctx, userID)
	if err != nil {
		return err
	}
	http.SetCookie(response, &http.Cookie{
		Name:     session.CookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	return nil
}

func clearSessionCookie(response http.ResponseWriter) {
	http.SetCookie(response, &http.Cookie{
		Name:     session.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

func sessionToken(request *http.Request) string {
	cookie, err := request.Cookie(session.CookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}

package health

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	customHttp "github.com/h-karlsson/uptime/internal/http"
)

const pingTimeout = 2 * time.Second

// GetHealth reports 503 when the database is unreachable.
func GetHealth(database Pinger) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		ctx, cancel := context.WithTimeout(request.Context(), pingTimeout)
		defer cancel()

		if err := database.Ping(ctx); err != nil {
			slog.Warn("health check failed", "error", err)
			customHttp.FormatResponseToJson(
				response,
				http.StatusServiceUnavailable,
				map[string]string{"status": "unavailable", "database": "down"},
			)
			return
		}
		customHttp.FormatResponseToJson(response, http.StatusOK, map[string]string{"status": "ok", "database": "up"})
	}
}

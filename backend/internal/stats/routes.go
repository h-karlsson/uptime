package stats

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func RegisterRoutes(router *http.ServeMux, stats *Service, requireUser func(http.Handler) http.Handler) {
	router.Handle("GET /api/monitors/{id}/stats", requireUser(customHttp.HandleJSON(GetSummary(stats))))
	router.Handle("GET /api/monitors/{id}/days", requireUser(customHttp.HandleJSON(ListDays(stats))))
	router.Handle("GET /api/monitors/{id}/incidents", requireUser(customHttp.HandleJSON(ListIncidents(stats))))
}

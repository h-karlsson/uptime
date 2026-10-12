package monitor

import (
	"net/http"

	customHttp "github.com/h-karlsson/uptime/internal/http"
)

func RegisterRoutes(router *http.ServeMux, monitors *Service, requireUser func(http.Handler) http.Handler) {
	router.Handle("GET /api/monitors", requireUser(customHttp.HandleJSON(List(monitors))))
	router.Handle("POST /api/monitors", requireUser(customHttp.HandleJSONStatus(http.StatusCreated, Create(monitors))))
	router.Handle("GET /api/monitors/{id}", requireUser(customHttp.HandleJSON(Get(monitors))))
	router.Handle("PUT /api/monitors/{id}", requireUser(customHttp.HandleJSON(Update(monitors))))
	router.Handle("DELETE /api/monitors/{id}", requireUser(customHttp.HandleNoContent(Delete(monitors))))
}

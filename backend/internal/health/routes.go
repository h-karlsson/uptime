package health

import "net/http"

func RegisterRoutes(router *http.ServeMux, database Pinger) {
	router.HandleFunc("GET /api/health", GetHealth(database))
}

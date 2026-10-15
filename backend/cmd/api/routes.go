package main

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
	"github.com/h-karlsson/uptime/internal/health"
	"github.com/h-karlsson/uptime/internal/monitor"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/stats"
	"github.com/h-karlsson/uptime/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(pool *pgxpool.Pool, sessions *session.Service) *http.ServeMux {
	users := user.NewService(user.NewPostgresStore(pool))
	monitors := monitor.NewService(monitor.NewPostgresStore(pool))
	monitorStats := stats.NewService(stats.NewPostgresStore(pool))
	requireUser := auth.RequireUser(sessions)

	mux := http.NewServeMux()
	health.RegisterRoutes(mux, pool)
	auth.RegisterRoutes(mux, users, sessions)
	monitor.RegisterRoutes(mux, monitors, requireUser)
	stats.RegisterRoutes(mux, monitorStats, requireUser)
	return mux
}

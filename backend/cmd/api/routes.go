package main

import (
	"net/http"

	"github.com/h-karlsson/uptime/internal/auth"
	"github.com/h-karlsson/uptime/internal/health"
	"github.com/h-karlsson/uptime/internal/session"
	"github.com/h-karlsson/uptime/internal/user"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterRoutes(pool *pgxpool.Pool) *http.ServeMux {
	users := user.NewService(user.NewPostgresStore(pool))
	sessions := session.NewService(session.NewPostgresStore(pool))

	mux := http.NewServeMux()
	health.RegisterRoutes(mux, pool)
	auth.RegisterRoutes(mux, users, sessions)
	return mux
}

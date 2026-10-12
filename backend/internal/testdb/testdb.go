// Package testdb gives integration tests an empty, migrated Postgres database. Tests that use
// it are skipped unless TEST_DATABASE_URL is set.
package testdb

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Packages run their tests in parallel against the same database; holding this advisory lock
// for the whole test makes them take turns.
const lockID = 727274

func Open(t *testing.T) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run the Postgres integration tests")
	}
	ctx := context.Background()

	pool, err := store.Open(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(pool.Close)
	if name := pool.Config().ConnConfig.Database; !strings.HasSuffix(name, "_test") {
		t.Fatalf("refusing to empty database %q: TEST_DATABASE_URL must point at a database ending in _test", name)
	}

	conn, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := conn.Exec(ctx, "SELECT pg_advisory_lock($1)", lockID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		conn.Exec(context.Background(), "SELECT pg_advisory_unlock($1)", lockID)
		conn.Release()
	})

	if err := store.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}
	Exec(t, pool, "TRUNCATE users RESTART IDENTITY CASCADE")
	if err := store.EnsureCheckPartitions(ctx, pool, time.Now()); err != nil {
		t.Fatal(err)
	}
	return pool
}

func Exec(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

// CreateUser inserts a user and returns its id.
func CreateUser(t *testing.T, pool *pgxpool.Pool, email string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		"INSERT INTO users (email, password_hash) VALUES ($1, 'unused') RETURNING id", email,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// CreateMonitor inserts an enabled monitor checking url and returns its id.
func CreateMonitor(t *testing.T, pool *pgxpool.Pool, userID int64, url string) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		"INSERT INTO monitors (user_id, name, url) VALUES ($1, $2, $2) RETURNING id", userID, url,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

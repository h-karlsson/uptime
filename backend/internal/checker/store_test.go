package checker

import (
	"context"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

// newTestStore connects to TEST_DATABASE_URL and empties it. The tests are skipped without it.
func newTestStore(t *testing.T) (*Store, *pgxpool.Pool) {
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
	if err := store.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}
	exec(t, pool, "TRUNCATE users, monitors, incidents, checks RESTART IDENTITY CASCADE")
	exec(t, pool, "INSERT INTO users (email, password_hash) VALUES ('checker@uptime.io', 'unused')")

	checkerStore := NewStore(pool)
	if err := checkerStore.EnsurePartitions(ctx, time.Now()); err != nil {
		t.Fatal(err)
	}
	return checkerStore, pool
}

func exec(t *testing.T, pool *pgxpool.Pool, sql string, args ...any) {
	t.Helper()
	if _, err := pool.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

func createMonitor(t *testing.T, pool *pgxpool.Pool, url string, enabled bool) int64 {
	t.Helper()
	var id int64
	err := pool.QueryRow(context.Background(),
		"INSERT INTO monitors (user_id, name, url, enabled) VALUES (1, $1, $1, $2) RETURNING id", url, enabled,
	).Scan(&id)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func claimedIDs(targets []target) []int64 {
	ids := make([]int64, len(targets))
	for i, target := range targets {
		ids[i] = target.monitorID
	}
	return ids
}

func TestClaimDueSkipsLockedMonitors(t *testing.T) {
	checkerStore, pool := newTestStore(t)
	ctx := context.Background()
	for _, url := range []string{"https://a.example", "https://b.example", "https://c.example"} {
		createMonitor(t, pool, url, true)
	}

	// Another instance claims one monitor in a transaction it hasn't committed yet.
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	locked, err := checkerStore.queries.WithTx(tx).ClaimDueMonitors(ctx, 1)
	if err != nil || len(locked) != 1 {
		t.Fatalf("claim in transaction: %v, %v", locked, err)
	}

	claimed, err := checkerStore.ClaimDue(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(claimed) != 2 || slices.Contains(claimedIDs(claimed), locked[0].ID) {
		t.Fatalf("claimed %v while monitor %d was locked; want the other two", claimedIDs(claimed), locked[0].ID)
	}

	// When the other instance rolls back, its monitor is due again; the rest are scheduled.
	if err := tx.Rollback(ctx); err != nil {
		t.Fatal(err)
	}
	claimed, err = checkerStore.ClaimDue(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(claimedIDs(claimed), []int64{locked[0].ID}) {
		t.Fatalf("after rollback claimed %v, want only %d", claimedIDs(claimed), locked[0].ID)
	}
}

func TestClaimDueSchedulesNextCheck(t *testing.T) {
	checkerStore, pool := newTestStore(t)
	ctx := context.Background()
	id := createMonitor(t, pool, "https://due.example", true)
	createMonitor(t, pool, "https://paused.example", false)

	claimed, err := checkerStore.ClaimDue(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(claimedIDs(claimed), []int64{id}) {
		t.Fatalf("claimed %v, want only the enabled monitor %d", claimedIDs(claimed), id)
	}

	var secondsUntilNext float64
	err = pool.QueryRow(ctx, "SELECT extract(epoch FROM next_check_at - now()) FROM monitors WHERE id = $1", id).Scan(&secondsUntilNext)
	if err != nil {
		t.Fatal(err)
	}
	if secondsUntilNext < 55 || secondsUntilNext > 60 {
		t.Fatalf("next check in %.0fs, want the 60s interval", secondsUntilNext)
	}
}

func TestRecordTracksStatusAndIncidents(t *testing.T) {
	checkerStore, pool := newTestStore(t)
	ctx := context.Background()
	id := createMonitor(t, pool, "https://flaky.example", true)

	record := func(ok bool) {
		t.Helper()
		result := result{checkedAt: time.Now(), ok: ok}
		if !ok {
			message := "timed out after 10s"
			result.err = &message
		}
		if err := checkerStore.Record(ctx, id, result); err != nil {
			t.Fatal(err)
		}
	}
	assertState := func(wantStatus string, wantOpen, wantResolved int) {
		t.Helper()
		var status string
		var open, resolved int
		err := pool.QueryRow(ctx, `
			SELECT m.current_status,
			       count(i.id) FILTER (WHERE i.resolved_at IS NULL),
			       count(i.id) FILTER (WHERE i.resolved_at IS NOT NULL)
			FROM monitors m LEFT JOIN incidents i ON i.monitor_id = m.id
			WHERE m.id = $1
			GROUP BY m.current_status`, id).Scan(&status, &open, &resolved)
		if err != nil {
			t.Fatal(err)
		}
		if status != wantStatus || open != wantOpen || resolved != wantResolved {
			t.Fatalf("status %s with %d open and %d resolved incidents; want %s, %d, %d",
				status, open, resolved, wantStatus, wantOpen, wantResolved)
		}
	}

	record(false)
	record(false)
	assertState("unknown", 0, 0)

	record(false)
	assertState("down", 1, 0)

	record(false)
	assertState("down", 1, 0)

	record(true)
	assertState("up", 0, 1)

	var checks int
	if err := pool.QueryRow(ctx, "SELECT count(*) FROM checks WHERE monitor_id = $1", id).Scan(&checks); err != nil {
		t.Fatal(err)
	}
	if checks != 5 {
		t.Fatalf("stored %d checks, want 5", checks)
	}
}

func TestRecordIgnoresDeletedMonitor(t *testing.T) {
	checkerStore, pool := newTestStore(t)
	id := createMonitor(t, pool, "https://gone.example", true)
	exec(t, pool, "DELETE FROM monitors WHERE id = $1", id)

	if err := checkerStore.Record(context.Background(), id, result{checkedAt: time.Now(), ok: true}); err != nil {
		t.Fatalf("recording a deleted monitor: %v", err)
	}
}

func TestEnsurePartitionsIsIdempotent(t *testing.T) {
	checkerStore, pool := newTestStore(t)
	ctx := context.Background()
	december := time.Date(2030, time.December, 15, 12, 0, 0, 0, time.UTC)

	for range 2 {
		if err := checkerStore.EnsurePartitions(ctx, december); err != nil {
			t.Fatal(err)
		}
	}
	var partitions []string
	rows, err := pool.Query(ctx, "SELECT relname FROM pg_class WHERE relname IN ('checks_2030_12', 'checks_2031_01') ORDER BY relname")
	if err != nil {
		t.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		partitions = append(partitions, name)
	}
	if !slices.Equal(partitions, []string{"checks_2030_12", "checks_2031_01"}) {
		t.Fatalf("partitions = %v, want December and the following January", partitions)
	}
}

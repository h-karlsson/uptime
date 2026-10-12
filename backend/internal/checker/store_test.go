package checker

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/testdb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newTestStore(t *testing.T) (*Store, *pgxpool.Pool, int64) {
	t.Helper()
	pool := testdb.Open(t)
	return NewStore(pool), pool, testdb.CreateUser(t, pool, "checker@uptime.io")
}

func claimedIDs(targets []target) []int64 {
	ids := make([]int64, len(targets))
	for i, target := range targets {
		ids[i] = target.monitorID
	}
	return ids
}

func TestClaimDueSkipsLockedMonitors(t *testing.T) {
	checkerStore, pool, userID := newTestStore(t)
	ctx := context.Background()
	for _, url := range []string{"https://a.example", "https://b.example", "https://c.example"} {
		testdb.CreateMonitor(t, pool, userID, url)
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
	checkerStore, pool, userID := newTestStore(t)
	ctx := context.Background()
	id := testdb.CreateMonitor(t, pool, userID, "https://due.example")
	paused := testdb.CreateMonitor(t, pool, userID, "https://paused.example")
	testdb.Exec(t, pool, "UPDATE monitors SET enabled = false WHERE id = $1", paused)

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
	checkerStore, pool, userID := newTestStore(t)
	ctx := context.Background()
	id := testdb.CreateMonitor(t, pool, userID, "https://flaky.example")

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
	checkerStore, pool, userID := newTestStore(t)
	id := testdb.CreateMonitor(t, pool, userID, "https://gone.example")
	testdb.Exec(t, pool, "DELETE FROM monitors WHERE id = $1", id)

	if err := checkerStore.Record(context.Background(), id, result{checkedAt: time.Now(), ok: true}); err != nil {
		t.Fatalf("recording a deleted monitor: %v", err)
	}
}

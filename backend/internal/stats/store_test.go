package stats

import (
	"context"
	"testing"
	"time"

	"github.com/h-karlsson/uptime/internal/store"
	"github.com/h-karlsson/uptime/internal/testdb"
	"github.com/jackc/pgx/v5/pgxpool"
)

func newIntegrationStore(t *testing.T) (*PostgresStore, *pgxpool.Pool, int64) {
	t.Helper()
	pool := testdb.Open(t)
	// The daily history reaches into the previous month.
	if err := store.EnsureCheckPartitions(context.Background(), pool, time.Now().AddDate(0, -1, 0)); err != nil {
		t.Fatal(err)
	}
	userID := testdb.CreateUser(t, pool, "stats@uptime.io")
	return NewPostgresStore(pool), pool, testdb.CreateMonitor(t, pool, userID, "https://stats.example")
}

func insertCheck(t *testing.T, pool *pgxpool.Pool, monitorID int64, checkedAt time.Time, ok bool, latencyMs *int) {
	t.Helper()
	testdb.Exec(t, pool,
		"INSERT INTO checks (monitor_id, checked_at, ok, latency_ms) VALUES ($1, $2, $3, $4)",
		monitorID, checkedAt, ok, latencyMs)
}

func ms(value int) *int {
	return &value
}

func TestSummaryAndLatency(t *testing.T) {
	statsStore, pool, monitorID := newIntegrationStore(t)
	ctx := context.Background()
	now := time.Now()
	latest := now.Add(-1 * time.Minute)

	insertCheck(t, pool, monitorID, now.Add(-4*time.Minute), true, ms(100))
	insertCheck(t, pool, monitorID, now.Add(-3*time.Minute), true, ms(200))
	insertCheck(t, pool, monitorID, now.Add(-2*time.Minute), false, nil)
	insertCheck(t, pool, monitorID, latest, true, ms(300))
	insertCheck(t, pool, monitorID, now.Add(-72*time.Hour), true, ms(1000))

	day, err := statsStore.Summary(ctx, monitorID, now.Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	// p95 of 100, 200 and 300 interpolates 90% of the way from 200 to 300.
	if day.Checks != 4 || *day.UptimePercent != 75 || *day.AvgLatencyMs != 200 || *day.P95LatencyMs != 290 {
		t.Fatalf("24h summary = %d checks, %v%% up, avg %v, p95 %v; want 4, 75, 200, 290",
			day.Checks, *day.UptimePercent, *day.AvgLatencyMs, *day.P95LatencyMs)
	}
	if day.LastCheckedAt == nil || day.LastCheckedAt.Sub(latest).Abs() > time.Millisecond {
		t.Fatalf("last checked at %v, want %v", day.LastCheckedAt, latest)
	}

	week, err := statsStore.Summary(ctx, monitorID, now.Add(-7*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if week.Checks != 5 || *week.UptimePercent != 80 || *week.AvgLatencyMs != 400 {
		t.Fatalf("7d summary = %d checks, %v%% up, avg %v; want 5, 80, 400", week.Checks, *week.UptimePercent, *week.AvgLatencyMs)
	}

	points, err := statsStore.Latency(ctx, monitorID, now.Add(-24*time.Hour), 15*time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	var checks, failed int
	for _, point := range points {
		checks += point.Checks
		failed += point.FailedChecks
		if !point.Time.Equal(point.Time.Truncate(15 * time.Minute)) {
			t.Fatalf("bucket %v doesn't start on a 15 minute boundary", point.Time)
		}
	}
	if checks != 4 || failed != 1 {
		t.Fatalf("latency buckets hold %d checks with %d failures, want 4 and 1", checks, failed)
	}
}

func TestSummaryWithoutChecks(t *testing.T) {
	statsStore, _, monitorID := newIntegrationStore(t)
	summary, err := statsStore.Summary(context.Background(), monitorID, time.Now().Add(-24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if summary.Checks != 0 || summary.UptimePercent != nil || summary.AvgLatencyMs != nil || summary.P95LatencyMs != nil || summary.LastCheckedAt != nil {
		t.Fatalf("summary = %+v, want zero checks and nil statistics", summary)
	}
}

func TestDays(t *testing.T) {
	statsStore, pool, monitorID := newIntegrationStore(t)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	twoDaysAgo := today.AddDate(0, 0, -2)

	insertCheck(t, pool, monitorID, twoDaysAgo.Add(10*time.Hour), true, ms(100))
	insertCheck(t, pool, monitorID, twoDaysAgo.Add(11*time.Hour), false, nil)
	insertCheck(t, pool, monitorID, today.Add(time.Second), true, ms(100))

	days, err := statsStore.Days(context.Background(), monitorID, today.AddDate(0, 0, -(DailyHistory-1)))
	if err != nil {
		t.Fatal(err)
	}
	want := []Day{
		{Date: twoDaysAgo.Format(time.DateOnly), Checks: 2, UptimePercent: 50},
		{Date: today.Format(time.DateOnly), Checks: 1, UptimePercent: 100},
	}
	if len(days) != len(want) || days[0] != want[0] || days[1] != want[1] {
		t.Fatalf("days = %+v, want %+v", days, want)
	}
}

func TestIncidentsNewestFirst(t *testing.T) {
	statsStore, pool, monitorID := newIntegrationStore(t)
	now := time.Now()
	testdb.Exec(t, pool, "INSERT INTO incidents (monitor_id, started_at, resolved_at, cause) VALUES ($1, $2, $3, 'timed out after 10s')",
		monitorID, now.Add(-2*time.Hour), now.Add(-time.Hour))
	testdb.Exec(t, pool, "INSERT INTO incidents (monitor_id, started_at, cause) VALUES ($1, $2, 'expected status 200, got 503')",
		monitorID, now.Add(-10*time.Minute))

	incidents, err := statsStore.Incidents(context.Background(), monitorID)
	if err != nil {
		t.Fatal(err)
	}
	if len(incidents) != 2 || incidents[0].ResolvedAt != nil || incidents[1].ResolvedAt == nil || *incidents[1].Cause != "timed out after 10s" {
		t.Fatalf("incidents = %+v, want the open one first, then the resolved one", incidents)
	}
}
